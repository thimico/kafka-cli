package kafka

import "github.com/Shopify/sarama"

func NewAdmin(addr []string, config *sarama.Config) (sarama.ClusterAdmin, error) {
	a, err := sarama.NewClusterAdmin(addr, config)
	return a, err
}
