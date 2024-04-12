
# db2db

A command line program to copy database tables and rows to another database.

Designed to replicate test environment db in local db to test issues locally.



## Usage

Download the appropriate executable file from [releases.](https://github.com/karataymarufemre/db2db/releases/tag/v1.0.0)


Edit the command below to connect to your own DBs.

```bash
  db2db -source sqlserver://<username>:<password>@<ip>:<port>?database=<databaseName> -target sqlserver://<username>:<password>@<ip>:<port>?database=<databaseName>
```

(optional)  run help to list all commands

```bash
  db2db -h
```

(optional)  exclude tables using -et or -excludeTables

```bash
  db2db -et table_1,table_2
```

(optional)  don't create tables using -ct -createTables

```bash
  db2db -ct=false
```
