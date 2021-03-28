package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/thimico/kafka-cli/kafka"
	"github.com/thimico/kafka-cli/log"
	"github.com/thimico/kafka-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

var consumerExample = `
# Consume msgs which start with given topic, partition and offset.
    ./kafka-cli consumer --bootstrapServers=localhost:9092 --topic=singed --offset=1 --partition=1
`

type consumerOptions struct {
	bootstrapServers string
	topic            string
	partition        int32
	offset           int64
}

func newConsumerOptions() *consumerOptions {
	return &consumerOptions{}
}

func (o *consumerOptions) run(cmd *cobra.Command, args []string) {
	if o.topic != "" {
		config := sarama.NewConfig()
		c, err := kafka.NewConsumer(strings.Split(o.bootstrapServers, ","), config)
		utils.CheckErr(err)
		defer func() {
			utils.CheckErr(c.Close())
		}()
		pc, err := c.ConsumePartition(o.topic, o.partition, o.offset)
		utils.CheckErr(err)
		defer func() {
			utils.CheckErr(pc.Close())
		}()
		for {
			select {
			case msg := <-pc.Messages():
				utils.PrintConsumerMessage(msg)
			case err := <-pc.Errors():
				log.Info("partition consumer", zap.Error(err))
			default:
				//do nothing
			}
		}
	} else {
		cmd.Help()
	}
}

func NewCmdConsumer() *cobra.Command {
	o := newConsumerOptions()
	cmd := &cobra.Command{
		Use:     "consumer",
		Short:   "Consume kafka message with given topic and partition",
		Long:    "Consume kafka message with given topic and partition, if you want to consume all partition, please refer to use consumerg",
		Example: consumerExample,
		Run:     o.run,
	}
	cmd.Flags().StringVarP(&o.bootstrapServers, "bootstrap-servers", "b", "localhost:9092", "The Kafka server to connect to.more than one should be separated by commas")
	cmd.Flags().StringVar(&o.topic, "topic", o.topic, "REQUIRED: The topics to consume,more than one should be separated by commas")
	cmd.Flags().Int32Var(&o.partition, "partition", 0, "The partition to consume (default 0)")
	cmd.Flags().Int64Var(&o.offset, "offset", sarama.OffsetNewest, "Which offset to consume start with, -2 means oldest, -1 means newest")
	return cmd
}
