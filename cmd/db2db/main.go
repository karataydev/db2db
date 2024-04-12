package main

import (
	"log"
	"slices"

	"github.com/karataymarufemre/db2db/internal/config"
	"github.com/karataymarufemre/db2db/internal/connection"
	"github.com/karataymarufemre/db2db/internal/data"
	"github.com/karataymarufemre/db2db/internal/table"
)

func main() {
	log.SetFlags(0)
	conf, err := config.FromFlagArgs()
	if err != nil {
		log.Fatal(err)
	}

	dbUrls := conf.DBs
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
	var columnDetailService table.ColumnDetailService
	var ddlGenerator table.DDLGenerator
	var creator table.Creator
	var reader data.Reader
	var writer data.Writer
	var replicator data.Replicator

	if dbUrls.Source.Driver == "mssql" {
		tableNameReader = table.NewSQLServerTableNameReader(db.Source)
		columnDetailService = table.NewSQLServerColumnDetailService(db.Source)
		ddlGenerator = table.NewSQLServerDDLGenerator(db.Source, columnDetailService)
		creator = table.NewSQLCreator(db.Target, &ddlGenerator)
		reader = data.NewSQLServerReader(db.Source, columnDetailService)
		writer = data.NewSQLServerWriter(db.Target)
		replicator = data.NewRowReplicator(reader, writer)
	} else {
		log.Fatal("Unsupported database driver")
	}

	// table names
	tables, err := tableNameReader.GetTableNames()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table Names:", tables)

	// generate ddl and apply it to target database
	if conf.CreateTables {
		if err = creator.CreateTables(tables); err != nil {
			log.Fatal(err)
		}
	}

	// row inserts
	for _, tableName := range tables {
		if slices.Contains(conf.ExcludedTables, tableName) {
			continue
		}
		err = replicator.Replicate(tableName)
		if err != nil {
			log.Fatal(err)
		}
	}

}
