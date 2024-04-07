package main

import (
	"log"

	"github.com/karataymarufemre/db2db/internal/config"
	"github.com/karataymarufemre/db2db/internal/connection"
	"github.com/karataymarufemre/db2db/internal/table"
)

func main() {
	log.SetFlags(0)
	dbUrls, err := config.DBUrlsFromFlagArgs()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connection string example: sqlserver://sa:password@localhost:1433?database=dbName")
	log.Printf("source connection: %+v\n", *dbUrls.Source)
	log.Printf("target connection: %+v\n", *dbUrls.Target)

	if dbUrls.Source.Driver != dbUrls.Target.Driver {
		log.Fatal("different source and target drivers not supported!")
	}

	db, err := connection.ConnectToDatabases(*dbUrls)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tableNameReader table.TableNameReader
	var ddlGenerator table.DDLGenerator

	if dbUrls.Source.Driver == "mssql" {
		tableNameReader = &table.SQLServerTableNameReader{}
		ddlGenerator = &table.SQLServerDDLGenerator{}
	} else {
		log.Fatal("Unsupported database driver")
	}

	// table names
	tables, err := tableNameReader.GetTableNames(db.Source)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table Names:", tables)

	// generate ddl and apply it to target database
	creator := table.NewSQLCreator(db, &ddlGenerator)
	if err := creator.Create(tables); err != nil {
		log.Fatal(err)
	}

	// row inserts

}
