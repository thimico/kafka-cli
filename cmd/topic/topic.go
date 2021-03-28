package topic

import (
	"github.com/Shopify/sarama"
	"github.com/thimico/kafka-cli/kafka"
	"github.com/thimico/kafka-cli/log"
	"github.com/thimico/kafka-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

var (
	topicExample = `
# Create a topic
    ./kafka-cli topic -c=singed  --partition-num=10 --replica-num=1
    result: 
        {"level":"info","ts":1612423941.894614,"caller":"log/log.go:16","msg":"Create topic success","topic":"singed","partition num":10,"replica num":1}

# List all available topics.
    ./kafka-cli topic -l

# List details for the given topics.more than one should be separated by commas
    ./kafka-cli topic --describe=singed
    result:
        ****************************************
		TOPIC:singed
		DETAIL:{
 		"Err": 0,
		"Name": "singed",
 		"IsInternal": false,
 		"Partitions": [
  		{
   			"Err": 0,
   			"ID": 0,
   			"Leader": 0,
			"Replicas": [
   				 0
   			],
   			"Isr": [
    			0
   			],
   			"OfflineReplicas": null
  		},
  		{
   			"Err": 0,
   			"ID": 1,
   			"Leader": 0,
   			"Replicas": [
    			0
   			],
   			"Isr": [
    			0
   			],
   			"OfflineReplicas": null
  		}
 		]
		}

# Delete a topic.
    ./kafka-cli topic -d=singed
    result:
		{"level":"info","ts":1612424432.454704,"caller":"log/log.go:16","msg":"Delete Topic success","topic":"singed"}

# Add partition number of topic
    ./kafka-cli topic --add-partition=singed --partition-num=3
    result:
		{"level":"info","ts":1612424575.056782,"caller":"log/log.go:16","msg":"Add partition success","topic":"singed","partition num":3}
`
)

type topicOptions struct {
	bootstrapServers string
	list             bool
	describe         string
	create           string
	delete           string
	addPartition     string
	numPartition     int32 //创建topic时指定的partition
	numReplica       int16 //创建topic时指定的副本数
}

func newTopicOptions() *topicOptions {
	return &topicOptions{}
}

func (o *topicOptions) run(cmd *cobra.Command, args []string) {
	config := sarama.NewConfig()
	servers := strings.Split(o.bootstrapServers, ",")
	admin, err := kafka.NewAdmin(servers, config)
	utils.CheckErr(err)
	defer func() {
		utils.CheckErr(admin.Close())
	}()
	if o.list {
		topics, err := admin.ListTopics()
		utils.CheckErr(err)
		for k, v := range topics {
			utils.PrintTopic(k, v)
		}
	} else if o.describe != "" {
		topics, err := admin.DescribeTopics(strings.Split(o.describe, ","))
		utils.CheckErr(err)
		for _, v := range topics {
			utils.PrintTopicMeta(v)
		}
	} else if o.create != "" {
		err := admin.CreateTopic(o.create, &sarama.TopicDetail{NumPartitions: o.numPartition, ReplicationFactor: o.numReplica}, false)
		utils.CheckErr(err)
		log.Info("Create topic success", zap.String("topic", o.create), zap.Int32("partition num", o.numPartition), zap.Int16("replica num", o.numReplica))
	} else if o.delete != "" {
		err := admin.DeleteTopic(o.delete)
		utils.CheckErr(err)
		log.Info("Delete Topic success", zap.String("topic", o.delete))
	} else if o.addPartition != "" {
		err := admin.CreatePartitions(o.addPartition, o.numPartition, [][]int32{}, false)
		utils.CheckErr(err)
		log.Info("Add partition success", zap.String("topic", o.addPartition), zap.Int32("partition num", o.numPartition))
	} else {
		cmd.Help()
	}
}

func NewCmdTopic() *cobra.Command {
	o := newTopicOptions()

	cmd := &cobra.Command{
		Use:     "topic",
		Short:   "Kafka topic operations",
		Long:    "Topic operations, include topic create、list、delete、detail, topic partition create",
		Example: topicExample,
		Run:     o.run,
	}
	cmd.Flags().StringVarP(&o.bootstrapServers, "bootstrap-server", "b", "localhost:9092", "The Kafka server to connect to.more than one should be separated by commas")
	cmd.Flags().BoolVarP(&o.list, "list", "l", o.list, "List all available topics.")
	cmd.Flags().StringVar(&o.describe, "describe", o.describe, "List details for the given topics.more than one should be separated by commas")
	cmd.Flags().StringVarP(&o.create, "create", "c", o.create, "Create a new topic.")
	cmd.Flags().Int32Var(&o.numPartition, "partition-num", 1, "The specified partition when create topic or add partition")
	cmd.Flags().Int16Var(&o.numReplica, "replica-num", 1, "The specified replica when create topic")
	cmd.Flags().StringVarP(&o.delete, "delete", "d", o.delete, "Delete a topic.")
	cmd.Flags().StringVar(&o.addPartition, "add-partition", o.addPartition, "The Topic which need to create partition, partition num must higher than which already exists")
	return cmd
}
