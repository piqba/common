# common packages

Introduction
This is a repository intended to serve common packages that can be used in our projects.


### Folder structure

```bash
├── Makefile
├── README.md
├── coverage.out
├── go.mod
├── go.sum
└── pkg
    ├── config
    │   └── config.go
    ├── databases
    │   ├── elasticsearch
    │   │   └── elasticsearch.go
    │   ├── mongo
    │   │   └── mongo.go
    │   ├── postgres
    │   │   └── postgres.go
    │   ├── redis
    │   │   └── redis.go
    │   └── types_options.go
    ├── httpsrv
    │   └── httpsrv.go
    ├── jwt
    │   ├── jwt.go
    │   └── jwt_test.go
    ├── logger
    │   └── logger.go
    ├── o11y
    │   ├── prometheus.go
    │   └── type_options.go
    └── tools
        ├── sets.go
        ├── time_checker.go
        └── time_checker_test.go

12 directories, 20 files

```

# How to use 

- Add this command
```bash
go env -w GOPRIVATE=github.com/piqba
```
- add into you ~/.gitconfig

```bash
[url "https://<USERNAME>:<TOKEN>@github.com/"]
insteadOf = https://github.com/
```

Then use `go get github.com/piqba/common`
