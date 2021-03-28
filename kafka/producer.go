package kafka

import (
	"github.com/Shopify/sarama"
)

func NewProducer(addrs []string, config *sarama.Config) (sarama.SyncProducer, error){
	return sarama.NewSyncProducer(addrs, config)
}
