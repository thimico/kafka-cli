package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/thimico/kafka-cli/kafka"
	"github.com/thimico/kafka-cli/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"strings"
)

var (
	consumergExample = `
# Consume kafka messages of set of topic which certain group id, and will print those message in terminal
	kafka-cli consume -t test,singed -c default -b localhost:9092
		`
)

type consumerGOptions struct {
	bootstrapServers string
	groupID          string
	topics           string
}

func newConsumerGOptions() *consumerGOptions {
	return &consumerGOptions{}
}

func (o *consumerGOptions) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (o *consumerGOptions) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (o *consumerGOptions) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		utils.PrintConsumerMessage(msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (o *consumerGOptions) run(cmd *cobra.Command, args []string) {
	if o.topics != "" {
		config := sarama.NewConfig()
		config.Consumer.Offsets.Initial = sarama.OffsetOldest

		c, err := kafka.NewConsumerGroup(strings.Split(o.bootstrapServers, ","), o.groupID, config)
		utils.CheckErr(err)
		defer func() {
			utils.CheckErr(c.Close())
		}()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = c.Consume(ctx, strings.Split(o.topics, ","), o)
		utils.CheckErr(err)
	} else {
		cmd.Help()
	}
}

func NewCmdConsumeGroup() *cobra.Command {
	o := newConsumerGOptions()
	cmd := &cobra.Command{
		Use:     "consumerg",
		Short:   "Consume kafka message with given topics and group_id",
		Long:    "Consume kafka message with given topics and group_id, and message will be auto committed",
		Example: consumergExample,
		Run:     o.run,
	}

	cmd.Flags().StringVarP(&o.bootstrapServers, "bootstrap-servers", "b", "localhost:9092", "The Kafka server to connect to.more than one should be separated by commas")
	cmd.Flags().StringVar(&o.topics, "topics", o.topics, "The topics to consume,more than one should be separated by commas")
	cmd.Flags().StringVar(&o.groupID, "group-id", "kafka-cli", "The consumer group ID")
	return cmd
}
