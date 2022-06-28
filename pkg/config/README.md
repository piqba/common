# Config module

This module, pkg or folder contain two methods `LoadConfig` & `LoadEnvOrFallback`

#### LoadEnvOrFallback

It's used for check & load `env`  variables.

#### LoadConfig

It`s used for load a config file that could be in different format like (json,yaml,toml, ...others)

For our projects the extension for our settings files was `yaml`

E.g:

```yaml
kraken:
  logLvl: INFO
  health_check:
    enable: true
    port: 19090
    path: "/metrics"
  env: dev
  repository:
    type: "postgresql" # (memory|postgresql|mongo)
    # postgresql
    pg_max_conn: "10"
    # If type is postgres this make  reference to a table
    # in other case like mongodb make reference to a collection
    db_collection_or_table: "info_kcore"
    pg_max_idle_conn: "10"
    pg_max_lifetime_conn: "60"
    pg_dsn_uri: "host=localhost port=5432 user=postgres password=password dbname=kcore sslmode=disable"
    # mongodbma
    mg_uri: "mongodb://localhost:27017/kcore"
    mg_db: "kcore"
  security:
    secretJWT: "jira-reports-developersecret"
  server:
    host: "0.0.0.0"
    port: 3000
    wrTimeOut: 60
    rdTimeOut: 60
    idleTimeOut: 60
```

E.g:

```go
func main() {
	cfg, err := config.LoadConfig(".", "settings", "yml")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(cfg.GetString("kraken.server.port"))
	//	APP_KRAKEN_SERVER_PORT=4000 go run main.go
	// result 4000
	//	go run main.go
	// result 3000
}

```