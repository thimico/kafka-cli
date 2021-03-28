package producer

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/thimico/kafka-cli/kafka"
	"github.com/thimico/kafka-cli/log"
	"github.com/thimico/kafka-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

var producerExample = `
# Produce a message
    ./kafka-cli producer --bootstrap-servers=localhost:9092 --key=13 --partitioner=random --topic=singed --value='test value'
    result:
        {"level":"info","ts":1612429377.79058,"caller":"log/log.go:16","msg":"Send message success","partition":2,"offset":0}
`

type producerOptions struct {
	bootstrapServers string
	topic            string
	key              string
	value            string
	headers          string
	partitioner      string
	partition        int32
}

func newProducerOptions() *producerOptions {
	return &producerOptions{}
}

func (o *producerOptions) validate() error {
	if o.topic == "" {
		return errors.New("empty topic")
	}
	if o.value == "" {
		return errors.New("empty value")
	}
	return nil
}

func (o *producerOptions) run(cmd *cobra.Command, args []string) {

	err := o.validate()
	if err != nil {
		log.Info("kafka producer validate flags failed", zap.String("reason", err.Error()))
		return
	}
	config := sarama.NewConfig()
	if o.partitioner == "random" {
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	} else if o.partition >= 0 {
		config.Producer.Partitioner = sarama.NewManualPartitioner
	}

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := kafka.NewProducer(strings.Split(o.bootstrapServers, ","), config)
	utils.CheckErr(err)
	defer func() {
		utils.CheckErr(producer.Close())
	}()

	msg := sarama.ProducerMessage{
		Topic: o.topic,
		Value: sarama.StringEncoder(o.value),
		//Metadata interface{}
		Partition: o.partition,
	}
	if o.key != "" {
		msg.Key = sarama.StringEncoder(o.key)
	}
	if o.headers != "" {
		var hdrs []sarama.RecordHeader
		arrHdrs := strings.Split(o.headers, ",")
		for _, h := range arrHdrs {
			if header := strings.Split(h, ":"); len(header) != 2 {
				utils.CheckErr(errors.New("-header should be key:value. Example: -headers=foo:bar,bar:foo"))
			} else {
				hdrs = append(hdrs, sarama.RecordHeader{
					Key:   []byte(header[0]),
					Value: []byte(header[1]),
				})
			}
		}
		if len(hdrs) != 0 {
			msg.Headers = hdrs
		}
	}
	partition, offset, err := producer.SendMessage(&msg)
	log.Info("Send message success", zap.Int32("partition", partition), zap.Int64("offset", offset))
}
func NewCmdProducer() *cobra.Command {
	o := newProducerOptions()
	cmd := &cobra.Command{
		Use:     "producer",
		Short:   "A kafka synchronous producer",
		Long:    "A kafka synchronous producer, with pretty much config options. but it's not asynchronous, which means it will wait for result before return",
		Example: producerExample,
		Run:     o.run,
	}
	cmd.Flags().StringVarP(&o.bootstrapServers, "bootstrap-servers", "b", "localhost:9092", "The Kafka server to connect to.more than one should be separated by commas")
	cmd.Flags().StringVar(&o.topic, "topic", o.topic, "REQUIRED: The topic id to produce messages to.")
	cmd.Flags().StringVar(&o.key, "key", "", "the key of message")
	cmd.Flags().StringVar(&o.value, "value", "", "REQUIRED: The message content which is going to be produced")
	cmd.Flags().StringVar(&o.partitioner, "partitioner", "hash", "The partitioning scheme to use. Can be hash, manual, or random")
	cmd.Flags().Int32Var(&o.partition, "partition", -1, "The partition which message produce to, if provided, it will use manual partitioner")
	cmd.Flags().StringVar(&o.headers, "headers", "", "The headers of the message. Example: -headers=foo:bar,bar:foo")
	return cmd
}
