# Database module

This module, pkg or folder contain all possibles `database` that can be used for our projects

- Postgresql (SQL)
- Redis (KV)
- MongoDB (NoSQL)
- ElasticSearch (Search Engine)
- BadgerDB (KV)


The file `types.go` contain useful structs for options definitions e.g:

```go

// PgOptions define some useful options for postgres db client
type PgOptions struct {
	MaxConnections         int
	MaxIdleConnections     int
	MaxLifeTimeConnections int
	// DbURI
	//  # Example DSN
	//  user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
	//
	//  # Example URL
	//  postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	DbURI string
}

// RdbOptions define some useful options for redis db client
type RdbOptions struct {
	// Addr
	//  # Example Addr
	//  "localhost:6379"
	Addr     string
	Username string
	Password string
	DB       int
}

```


