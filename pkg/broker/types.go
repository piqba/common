package broker

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
