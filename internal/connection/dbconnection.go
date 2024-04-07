package connection

import (
	"database/sql"
	"log"

	// drivers
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/karataymarufemre/db2db/internal/config"
)

type DBConnection struct {
	Source *sql.DB
	Target *sql.DB
}

func ConnectToDatabases(urls config.DBUrls) (*DBConnection, error) {
	source, err := getDB(urls.Source.Driver, urls.Source.ConnectionString)
	if err != nil {
		return nil, err
	}
	log.Println("Source Connected!")
	target, err := getDB(urls.Target.Driver, urls.Target.ConnectionString)
	if err != nil {
		return nil, err
	}
	log.Println("Target Connected!")
	return &DBConnection{Source: source, Target: target}, nil
}

func getDB(driver string, connectionString string) (*sql.DB, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (conn *DBConnection) Close() {
	err := conn.Target.Close()
	err = conn.Source.Close()
	if err != nil {
		log.Fatal(err)
	}
}
