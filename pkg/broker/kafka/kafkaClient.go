package kfk

import (
	"context"
	"github.com/piqba/common/pkg/broker"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//Used to execute client creation procedure only once.
var kafkaOnce sync.Once

// ProducerKafkaOnce return an instance of *kafka.Producer
func ProducerKafkaOnce(options broker.KafkaOptions) (*kafka.Producer, error) {
	var errK error
	var producerInstance *kafka.Producer
	kafkaOnce.Do(func() {
		p, err := kafka.NewProducer(
			&kafka.ConfigMap{
				"bootstrap.servers":        options.Host,
				"compression.type":         options.CompressionType, // better "snappy"
				"acks":                     options.Acks,            // all acks -1
				"enable.idempotence":       options.Idempotence,     // true
				"max.in.flight":            options.MaxInFlight,     // 5
				"message.send.max.retries": options.MaxRetries,      // 50
				"linger.ms":                options.MS,              // 25
			},
		)
		if err != nil {
			errK = err
		}
		producerInstance = p
	})
	return producerInstance, errK
}

// ProducerKafka return an instance of *kafka.Producer
func ProducerKafka(options broker.KafkaOptions) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers":        options.Host,
			"compression.type":         options.CompressionType, // better "snappy"
			"acks":                     options.Acks,            // all acks -1
			"enable.idempotence":       options.Idempotence,     // true
			"max.in.flight":            options.MaxInFlight,     // 5
			"message.send.max.retries": options.MaxRetries,      // 50
			"linger.ms":                options.MS,              // 25
		},
	)
	if err != nil {
		return nil, err
	}
	return p, err
}

// ConsumerKafka return an instance of *kafka.Consumer
// This is can be use to produce events.
func ConsumerKafka(options broker.KafkaOptions) (*kafka.Consumer, error) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               options.Host,
		"group.id":                        options.GroupID,
		"session.timeout.ms":              options.TimeOUT, // 6000
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		// Enable generation of PartitionEOF when the
		// end of a partition is reached.
		"enable.partition.eof": true,
		"auto.offset.reset":    options.OffsetType, // latest | earliest
		// handling kafka.Stats events (see below).
		"statistics.interval.ms": options.Statics, // 5000
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// CreateTopicIfNotExists check if topic exist ...
func CreateTopicIfNotExists(ctx context.Context, consumer *kafka.Consumer, topics []string, options broker.KafkaOptions) error {
	adminClient, err := kafka.NewAdminClientFromConsumer(consumer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		return err
	}

	for _, topic := range topics {
		_, err := adminClient.CreateTopics(
			ctx,

			[]kafka.TopicSpecification{
				{
					Topic:             topic,
					NumPartitions:     options.NumPartitions,     // 1
					ReplicationFactor: options.ReplicationFactor, // 1
				},
			},

			kafka.SetAdminOperationTimeout(maxDur))
		if err != nil {
			return err
		}
	}
	adminClient.Close()
	return nil
}
