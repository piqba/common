# Broker module

This module, pkg or folder contain two possible datasource for now only soppourt `kafka` but in the future could be `rabbitmq` & `redisStream`


In the `types.go` are all possible options
that can be used on the process for producer, consumer & partitions creation.

```go
// KafkaOptions ...
type KafkaOptions struct {
	Host              string
	CompressionType   string
	Acks              int
	Idempotence       bool
	MaxInFlight       int
	MaxRetries        int
	MS                int
	GroupID           string
	TimeOUT           int
	OffsetType        string
	Statics           int
	NumPartitions     int
	ReplicationFactor int
}

```


### TODO
- Implementar estrategia de broker