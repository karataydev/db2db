package config

import (
	"flag"
	"fmt"
)

type DBConfig struct {
	ConnectionString string
	Driver           string
}

type DBUrls struct {
	Source *DBConfig
	Target *DBConfig
}

func DBUrlsFromFlagArgs() (*DBUrls, error) {
	dbUrls := &DBUrls{Source: &DBConfig{}, Target: &DBConfig{}}
	flag.StringVar(&dbUrls.Source.ConnectionString, "source", "", "db connection string to copy from example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Source.ConnectionString, "s", "", "db connection string to copy from example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Source.Driver, "sourceDriver", "mssql", "db sql driver, available: \"mssql\", others will supported later.")

	flag.StringVar(&dbUrls.Target.ConnectionString, "target", "", "db connection string to copy to example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Target.ConnectionString, "t", "", "db connection string to copy to example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Target.Driver, "fromDriver", "mssql", "db sql driver, available: \"mssql\", others will supported later.")

	flag.Parse()

	if dbUrls.Source.ConnectionString == "" {
		return nil, fmt.Errorf("source flag is required use -source \"a connection string\"")
	}

	if dbUrls.Target.ConnectionString == "" {
		return nil, fmt.Errorf("target flag is required use -target \"a connection string\"")
	}

	return dbUrls, nil
}
