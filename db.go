package main

import (
	"database/sql"
	"fmt"
	"github.com/dimau/hh-api-client-go"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func initializePostgresConnection(postgresUserName, postgresPasswd, postgresServerName, postgresPort, postgresDB string) *sql.DB {
	connectionURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", postgresUserName, postgresPasswd, postgresServerName, postgresPort, postgresDB)
	db, err := sql.Open("pgx", connectionURL)
	failOnError(err, "Failed to create a connection to Postgres")
	return db
}

func insertVacancy(db *sql.DB, vacancy *hh.Vacancy) {

	stmt, err := db.Prepare(`INSERT INTO public."Vacancies" (title, description, link) VALUES ($1, $2, $3);`)
	failOnError(err, "Error when preparing a statement for inserting new vacancies to DB")

	_, err = stmt.Exec(vacancy.Name, vacancy.Snippet.Requirement, vacancy.Url)
	failOnError(err, "Error when inserting a new vacancy to DB")
}
