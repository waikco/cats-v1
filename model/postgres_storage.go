package model

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	json "github.com/json-iterator/go"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/waikco/cats-v1/conf"
)

//CreateTableQuery is sql query for creating fda_data table
const CreateTableQuery string = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS cats (
id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
name TEXT NOT NULL,
color TEXT NOT NULL,
age INT NOT NULL
);`

//TestCreateTableQuery is sql query for creating fda_data table
const TestCreateTableQuery string = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS cats (
id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
name TEXT NOT NULL,
color TEXT NOT NULL,
age INT NOT NULL
);`

type PostGres struct {
	database *sqlx.DB
	dbName   string
}

func BootstrapPostgres(config conf.Database) (Storage, error) {
	// connect to database
	dbInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName)
	db, err := sqlx.Connect("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	// return db connection
	storage := &PostGres{db, config.DatabaseName}
	_, err = storage.database.Exec(CreateTableQuery)
	if err != nil {
		return storage, err
	} else {
		log.Debug().Msg("table presence confirmed")
	}
	return storage, nil
}

func (p *PostGres) Insert(b []byte) (string, error) {
	var cat Cat
	if err := json.Unmarshal(b, &cat); err != nil {
		return "", err
	}

	var id string
	query := `INSERT INTO cats (name,color,age) VALUES ($1,$2, $3) RETURNING id`
	err := p.database.QueryRow(query, cat.Name, cat.Color, cat.Age).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (p *PostGres) Select(id string) ([]byte, error) {
	var cat Cat
	query := `SELECT name, color, age FROM cats WHERE id=$1`
	err := p.database.QueryRow(query, id).Scan(&cat.Name, &cat.Color, &cat.Age)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cat)
}

func (p *PostGres) SelectAll(count int, start int) ([]byte, error) {
	rows, err := p.database.Query(`SELECT id, name, color, age FROM cats LIMIT $1 OFFSET $2`, count, start)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var cats []Cat
	for rows.Next() {
		var cat Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Color, &cat.Age)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	if len(cats) == 0 {
		return nil, sql.ErrNoRows
	}
	return json.Marshal(cats)
}

func (p *PostGres) Update(id string, b []byte) error {
	var cat Cat
	if err := json.Unmarshal(b, &cat); err != nil {
		return err
	}

	query := `UPDATE cats SET name=$1, color=$2, age=$3 WHERE id=$4`
	_, err := p.database.Exec(query, cat.Name, cat.Color, cat.Age, id)
	return err
}

func (p *PostGres) Delete(id string) error {
	_, err := p.database.Exec("DELETE FROM cats where id=$1", id)
	return err
}

func (p *PostGres) Status() error {
	err := p.database.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostGres) Purge(table string) error {
	if _, err := p.database.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
		return fmt.Errorf("Error purging %s table: %v", table, err)
	}
	log.Info().Msgf("Purging %s table", table)
	if _, err := p.database.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)); err != nil {
		return err
	}
	return nil
}
