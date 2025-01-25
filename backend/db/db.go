package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
create table if not exists "user"
(
    id            uuid primary key default gen_random_uuid(),
    username      text unique not null,
    password      text        not null
);`

type DbConfig struct {
	name     string
	username string
	host     string
	password string
	port     int
}

func readConfigEnv() DbConfig {
	database_name, ok := os.LookupEnv("DB_DATABASE")
	if !ok {
		log.Fatal("Not database name, cannot start")
	}

	db_username, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		log.Fatal("Not database name, cannot start")
	}

	db_host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatal("Not database name, cannot start")
	}

	db_password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("Not database name, cannot start")
	}

	return DbConfig{
		name:     database_name,
		username: db_username,
		host:     db_host,
		password: db_password,
		port:     5432,
	}
}

func NewDb() *sqlx.DB {
	config := readConfigEnv()

	
	db, err := connectDB(config)
	if err != nil {
		//TODO: retry connect 
		log.Fatal(err)
	}

	db.MustExec(schema)

	return db
}

func connectDB(config DbConfig) (*sqlx.DB, error) {
	db_string := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.host, config.port, config.username, config.name, config.password)
	db, err := sqlx.Connect("postgres", db_string)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func defaultConfig() DbConfig {
	return DbConfig{
		name:     "janko",
		username: "janko",
		host:     "localhost",
		password: "janko",
		port:     5432,
	}
}

func NewTestDb() (*sqlx.DB, error) {
	config := defaultConfig()

	db, err := connectDB(config)
	if err != nil {
		return nil, err
	}

	test_db_name := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	create_str := fmt.Sprintf("create database \"%s\" ", test_db_name)
	_, err = db.Exec(create_str)
	if err != nil {
		return nil, err
	}

	config.name = test_db_name
	db, err = connectDB(config)
	if err != nil {
		return nil, err
	}

	db.MustExec(schema)

	return db, nil
}
