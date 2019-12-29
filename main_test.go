package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/waikco/cats-v1/cmd/server"
	"github.com/waikco/cats-v1/conf"
)

var a server.App
var config conf.Config

func TestMain(m *testing.M) {

	if os.Getenv("CONFIG_SWITCH") == "local" {
		config = conf.SaneDefaults()
	}
	viperInstance := viper.GetViper()
	err := viperInstance.UnmarshalExact(&config)
	if err != nil {
		log.Panic().Msgf("error parsing config: %s", err.Error())
	}
	a.Config = config
	a.Bootstrap()
	log.Info().Msgf("database ssl status is %s", a.Config.Database.SslMode)

	confirmTableexists()

	code := m.Run()

	purgeTable()

	os.Exit(code)

}

func purgeTable() {

}

func confirmTableexists() {

}

func TestEmptyTable(t *testing.T) {
	purgeTable()

	req, _ := http.NewRequest("GET", "/cats/v1/cats", nil)
	_ = executeRequest(req)

	//checkResponseCode(t, http.StatusOK, response.Code)
	//checkResponseBody(t, "[]", response.Body.String())
}

func TestGetNonExistentCat(t *testing.T) {
	purgeTable()

	req, _ := http.NewRequest("GET", "/cats/v1/13", nil)
	_ = executeRequest(req)

	//checkResponseCode(t, httptp.StatusOK, response.Code)
	//checkResponseBody(t, "[]", response.Body.String())
}

func TestCreateCat(t *testing.T) {
	purgeTable()

	//m :=
	//if m["id"] != 1.0 {
	//	t.Errorf("Expected: '1' | Received: '%v'", m["id"])
	//}
	//
	//if m["name"] != "test pet" {
	//	t.Errorf("Expected: 'test pet' | Received: %v", m["name"])
	//}
	//
	//if m["kind"] != "dog" {
	//	t.Errorf("Expected: 'dog' | Received: %v", m["kind"])
	//}
	//
	//if m["color"] != "red" {
	//	t.Errorf("Expected: 'red' | Received: %v", m["color"])
	//}
	//
	//if m["age"] != 3.0 {
	//	t.Errorf("Expected: 3 | Received: %v", m["age"])
	//}
}

func TestMassCreateCat(t *testing.T) {
	purgeTable()
	cargo := []byte(`[{"name":"test pet","kind":"dog","color":"red","age":0}, {"name":"test pet","kind":"dog","color":"red","age":1}]`)
	req, _ := http.NewRequest("POST", "/cats/v1/bulkpetadd", bytes.NewBuffer(cargo))
	response := executeRequest(req)

	//checkResponseCode(t, http.StatusCreated, response.Code)

	var n []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &n)

	for i, v := range n {

		if v["id"] != i {
			t.Errorf("Expected: '1' | Received: '%v'", v["id"])
		}

		if v["name"] != "test pet" {
			t.Errorf("Expected: 'test pet' | Received: %v", v["name"])
		}

		if v["kind"] != "dog" {
			t.Errorf("Expected: 'dog' | Received: %v", v["kind"])
		}

		if v["color"] != "red" {
			t.Errorf("Expected: 'red' | Received: %v", v["color"])
		}
		if v["age"] != i {
			t.Errorf("Expected: %v | Received: %v", i, v["age"])
		}

	}

}

func TestGetCat(t *testing.T) {
	purgeTable()
	addCats(1)

	_ = httptest.NewRequest("GET", "/cats/v1/1", nil)

	//response := executeRequest(req)

	//checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateCat(t *testing.T) {
	purgeTable()
	addCats(1)

	req, err := http.NewRequest("GET", "/cats/v1/1", nil)
	if err != nil {
		log.Panic().Msgf("error creating request: %s", err.Error())
	}
	response := executeRequest(req)
	//checkResponseCode(t, http.StatusOK, response.Code)

	var originalCat map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCat)

	cargo := []byte(`{"name":"updated pet name","kind":"updated dog","color":"new red","age":5}`)

	req, err = http.NewRequest("PUT", "/cats/v1/1", bytes.NewBuffer(cargo))
	if err != nil {
		log.Panic().Msgf("error creating request: %s", err.Error())
	}
	response = executeRequest(req)

	//checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	//
	//if m["id"] != originalCat["id"] {
	//	t.Errorf("Expected id to remain the same (%v). Got %v", originalCat["id"], m["id"])
	//}
	//
	//if m["name"] == originalCat["name"] {
	//	t.Errorf("Expected the name to change from '%v' to '%v", originalCat["name"], m["name"])
	//}
	//
	//if m["kind"] == originalCat["kind"] {
	//	t.Errorf("Expected the kind to change from '%v' to '%v", originalCat["kind"], m["kind"])
	//
	//}
	//if m["color"] == originalCat["color"] {
	//	t.Errorf("Expected the color to change from '%v' to '%v", originalCat["color"], m["color"])
	//}
	//
	//if m["age"] == originalCat["age"] {
	//	t.Errorf("Expected the age to change from '%v' to '%v", originalCat["age"], m["age"])
	//}
}

func TestDeleteCat(t *testing.T) {
	purgeTable()
	addCats(1)

	req, _ := http.NewRequest("GET", "/cats/v1/1", nil)
	_ = executeRequest(req)
	//checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/cats/v1/1", nil)
	_ = executeRequest(req)
	//checkResponseCode(t, http.StatusOK, response.Code)
	//checkResponseBody(t, `{"result":"success"}`, response.Body.String())

	req, _ = http.NewRequest("GET", "/cats/v1/1", nil)
	_ = executeRequest(req)
	//checkResponseCode(t, http.StatusOK, response.Code)
	//checkResponseBody(t, `[]`, response.Body.String())

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

//ConfirmTableexists checks for existence of app database
// todo move this into functional test
//func confirmTableexists() {
//	/*
//		if _, err := a.DB.Exec(createTableQuery); err != nil {
//			log.Fatal(err)
//		}
//	*/
//
//	log.Print("confirming table exists...")
//	if _, err := a.Storage.Exec(model.SqlTestCatsCreateTableQuery); err != nil {
//		log.Fatal().Msgf("Error creating cats table: " + err.Error())
//	}
//	log.Print("database confirmed...")
//
//}

//PurgeTable deletes items from table
// todo Move this into functional test
//func purgeTable() {
//	if _, err := a.AppDatabase.Exec("DELETE FROM cats"); err != nil {
//		log.Fatal().Msgf("Error purging cats table: " + err.Error())
//	}
//	log.Print("Purging cats table")
//	a.AppDatabase.Exec("ALTER SEQUENCE cats_id_seq RESTART WITH 1")
//}

//AddCats adds a variable number of records to the database
func addCats(quant int) {
	if quant < 1 {
		quant = 1
	}

	//for i := 0; i < quant; i++ {
	//	_, err := a.Storage.Insert(model.Cat{
	//		ID:    uuid.NewV4().String(),
	//		Name:  "cat" + strconv.Itoa(i),
	//		Color: "color" + strconv.Itoa(i),
	//		Age:   i,
	//	})
	//	if err != nil {
	//		log.Fatal().Err(err)
	//	}
	//}
}
