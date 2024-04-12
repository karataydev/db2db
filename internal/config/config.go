package config

import (
	"flag"
	"fmt"
	"strings"
)

type DBConfig struct {
	ConnectionString string
	Driver           string
}

type DBUrls struct {
	Source *DBConfig
	Target *DBConfig
}

type Config struct {
	DBs            *DBUrls
	CreateTables   bool
	ExcludedTables []string
}

type CommaSeparatedStringArray []string

func (c *CommaSeparatedStringArray) String() string {
	return strings.Join(*c, ",")
}

func (c *CommaSeparatedStringArray) Set(value string) error {
	*c = append(*c, strings.Split(value, ",")...)
	return nil
}

func FromFlagArgs() (*Config, error) {
	dbUrls := &DBUrls{Source: &DBConfig{}, Target: &DBConfig{}}
	conf := &Config{DBs: dbUrls, CreateTables: true}
	flag.StringVar(&dbUrls.Source.ConnectionString, "source", "", "db connection string to copy from example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Source.ConnectionString, "s", "", "db connection string to copy from example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Source.Driver, "sourceDriver", "mssql", "db sql driver, available: \"mssql\", others will supported later.")

	flag.StringVar(&dbUrls.Target.ConnectionString, "target", "", "db connection string to copy to example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Target.ConnectionString, "t", "", "db connection string to copy to example: \"sqlserver://sa:password@localhost:1433?database=dbName\"")
	flag.StringVar(&dbUrls.Target.Driver, "fromDriver", "mssql", "db sql driver, available: \"mssql\", others will supported later.")

	flag.BoolVar(&conf.CreateTables, "createTables", true, "create tables")
	flag.BoolVar(&conf.CreateTables, "ct", true, "create tables")

	var excludedTables CommaSeparatedStringArray
	flag.Var(&excludedTables, "excludeTables", "A list of comma-separated excluded table names to insert")
	flag.Var(&excludedTables, "et", "A list of comma-separated excluded table names to insert")

	flag.Parse()
	conf.ExcludedTables = excludedTables

	if dbUrls.Source.ConnectionString == "" {
		return nil, fmt.Errorf("source flag is required use -source \"a connection string\"")
	}

	if dbUrls.Target.ConnectionString == "" {
		return nil, fmt.Errorf("target flag is required use -target \"a connection string\"")
	}

	return conf, nil
}
