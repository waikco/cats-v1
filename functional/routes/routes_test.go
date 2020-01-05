package routes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/jmoiron/sqlx"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/waikco/cats-v1/functional"
	"github.com/waikco/cats-v1/model"
)

const urlStart = "http://localhost:8080"

var db *sqlx.DB

func TestMain(m *testing.M) {
	process, output, err := functional.StartBinary("config.yaml")
	if err != nil {
		log.Fatal().Msgf(
			"starting binary failed with error %v and out %s",
			err,
			output.String())
	}

	if err := confirmTable(); err != nil {
		log.Fatal().Err(err)
	}

	if err := addCats(5); err != nil {
		log.Fatal().Err(err)
	}

	defer func() { purgeTable() }()

	code := m.Run()

	if process != nil {
		_ = process.Kill()
	}

	os.Exit(code)

}

func purgeTable() {
	_, _ = db.Exec(`DELETE FROM cats`)
}

func TestEmptyTable(t *testing.T) {
	purgeTable()

	req, err := http.NewRequest("GET", urlStart+"/cats/v1/cats", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("err making request: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("unexpcted status code: got %d, want %d", response.StatusCode, http.StatusOK)
	}
	if body, err := ioutil.ReadAll(response.Body); err != nil && string(body) != `[]` {
		t.Errorf("unexpcted body: got %s, want %s", string(body), `[]`)
	}
}

func TestGetNonExistentCat(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, urlStart+"/cats/v1/cats/9C679C16-C38B-48A3-A645-A6F2457A49BC", nil)
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		t.Errorf("err making request: %v", err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: got %d, want %d", response.StatusCode, http.StatusNotFound)
	}
	if body, err := ioutil.ReadAll(response.Body); err != nil && string(body) != `[]` {
		t.Errorf("unexpected body: got %s, want %s", string(body), `[]`)
	}
}

func TestCRUD(t *testing.T) {
	payload := `{"name":"cat-1","color":"color-1","age":1}`
	var id uuid.UUID

	t.Run("create cat", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost,
			urlStart+"/cats/v1/",
			bytes.NewBuffer([]byte(payload)))
		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			t.Errorf("err making request: %v", err)
		}

		if response.StatusCode != http.StatusCreated {
			t.Errorf("unexpected status code: got %d, want %d", response.StatusCode, http.StatusCreated)
		}

		var m map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&m); err != nil {
			t.Fatalf("error reading bod: %v", err)
		}
		id, _ = uuid.FromString(m["result"].(string))
		if id == uuid.Nil {
			t.Errorf("unexpected result: got %s, want valid UUID", m["result"].(string))
		}
		if errs := m["errors"]; errs != nil {
			t.Errorf("unexpected errors in response: %v", m["errors"])
		}
	})

	t.Run("get cat", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, urlStart+"/cats/v1/cats/"+id.String(), nil)
		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			t.Errorf("err making request: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpcted status code: got %d, want %d", response.StatusCode, http.StatusOK)
		}
		if body, _ := ioutil.ReadAll(response.Body); string(body) != payload {
			t.Errorf("unexpected body: got %s, want %s", string(body), payload)
		}
	})

	var newCat = `{"name":"new-cat-1","color":"orange","age":3}`
	t.Run("update cat", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, urlStart+"/cats/v1/"+id.String(),
			bytes.NewBuffer([]byte(newCat)))
		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			t.Errorf("err making request: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpcted status code: got %d, want %d", response.StatusCode, http.StatusOK)
		}
		var m map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&m); err != nil {
			t.Fatalf("error reading body: %v", err)
		}
		if m["result"] != "success" {
			t.Errorf("unexpected result: got %+v, want `result` to equal `success`", m)
		}

		req, _ = http.NewRequest(http.MethodGet, urlStart+"/cats/v1/cats/"+id.String(), nil)
		response, err = client.Do(req)
		if err != nil {
			t.Errorf("err making request: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: got %d, want %d", response.StatusCode, http.StatusOK)
		}
		if body, _ := ioutil.ReadAll(response.Body); string(body) != newCat {
			t.Errorf("unexpected new cat: got %s, want %s", string(body), newCat)
		}
	})

	t.Run("delete cat", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, urlStart+"/cats/v1/cats/"+id.String(), nil)
		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			t.Errorf("err making request: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpcted status code: got %d, want %d", response.StatusCode, http.StatusOK)
		}
		var m map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&m); err != nil {
			t.Fatalf("error reading body: %v", err)
		}
		if m["result"] != "success" {
			t.Errorf("unexpected result: got %+v, want `result` to equal `success`", m)
		}
	})
}

// confirmTable ensures connectivity to the database, and existence of required table.
func confirmTable() error {
	dbInfo :=
		"host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	d, err := sqlx.Connect("postgres", dbInfo)
	if err != nil {
		return errors.Wrap(err, "Error connecting to database: %v")
	}
	db = d

	if _, err := db.Exec(model.TestCreateTableQuery); err != nil {
		return errors.Wrap(err, "error creating table")
	}
	return nil
}

// addCats adds a variable number of records to the database.
func addCats(quant int) error {
	if quant < 1 {
		quant = 1
	}
	log.Info().Msgf("adding %d rows of dummy data", quant)
	query := `INSERT INTO cats (name,color,age) VALUES ($1,$2, $3) RETURNING id`
	for i := 0; i < quant; i++ {
		_, err := db.Exec(query,
			fmt.Sprintf("cat-%v", i),
			fmt.Sprintf("color-%v", i),
			i)

		if err != nil {
			return err
		}
	}
	return nil
}
