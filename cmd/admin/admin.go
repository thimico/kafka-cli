package admin

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/thimico/kafka-cli/kafka"
	"github.com/thimico/kafka-cli/log"
	"github.com/thimico/kafka-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

/*
describeCluster
DescribeLogDirs
 */

var adminExample = `
# List all consumer groups
    ./kafka-cli admin --list-consumer-groups

# Describe a certain consumer group
    ./kafka-cli admin --describe-groups --groups=garvin

# Delete consumer groups
    ./kafka-cli admin --delete-groups --groups=garvin

# Delete message records, records which less than specified offset will be deleted
    ./kafka-cli admin --delete-records --topics=singed --partitions=1,2 --offset=100

# List consumer offset
    ./kafka-cli admin --list-consumer-offsets --groups=garvin --topics=singed,test --partitions=0,1,2

# Describe cluster
    ./kafka-cli admin --describe-cluster

# Describe log dirs
    ./kafka-cli admin --describe-log-dirs --brokers=0,1   
`

type adminOptions struct {
	bootstrapServers string
	groups string
	topics string
	partitions string
	brokers string
	offset int64

	deleteRecords bool
	listConsumerGroups bool
	describeGroups bool
	deleteGroups bool
	listConsumerOffsets bool
	describeCluster bool
	describeLogDirs bool
}

func newAdminOptions() *adminOptions {
	return &adminOptions{}
}

func (o *adminOptions) validate() error {
	if o.deleteRecords {
		if o.topics == "" || o.partitions == "" || o.offset == 0 {
			return errors.New("when delete records, topics or partitions and offset should not be empty")
		}
	}
	if o.describeGroups {
		if o.groups == "" {
			return errors.New("when describe groups, groups flags should not be empty")
		}
	}
	if o.deleteGroups {
		if o.groups == "" {
			return errors.New("when delete groups, groups flag should not be empty")
		}
	}
	if o.listConsumerOffsets {
		if o.topics == "" || o.partitions == "" {
			return errors.New("when list consumer offsets, groups, topics and partitions should not be empty")
		}
	}
	if o.describeLogDirs {
		if o.brokers == "" {
			return errors.New("when list consumer offsets, groups, topics and partitions should not be empty")
		}
	}
	return nil
}

func (o *adminOptions) run(cmd *cobra.Command, args []string) {
	err := o.validate()
	if err != nil {
		log.Info("admin flags validate failed", zap.Error(err))
		return
	}
	config := sarama.NewConfig()
	servers := strings.Split(o.bootstrapServers, ",")
	admin, err := kafka.NewAdmin(servers, config)
	utils.CheckErr(err)
	defer func() {
		utils.CheckErr(admin.Close())
	}()

	if  o.deleteRecords {
		topics := strings.Split(o.topics, ",")
		_partitions := strings.Split(o.partitions, ",")
		partitionOffsets := map[int32]int64{}
		for _, p := range _partitions {
			partition, err := strconv.ParseInt(p, 10, 64)
			utils.CheckErr(err)
			partitionOffsets[int32(partition)]=o.offset
		}
		for _, t := range topics {
			utils.CheckErr(admin.DeleteRecords(t, partitionOffsets))
			log.Info("delete records success", zap.String("topic", t), zap.Any("offsetPartitions", partitionOffsets))
		}
	} else if o.listConsumerGroups{
		groups, err := admin.ListConsumerGroups()
		utils.CheckErr(err)
		utils.PrintConsumerGroups(groups)
	} else if o.describeGroups {
		groupDetail, err := admin.DescribeConsumerGroups(strings.Split(o.groups, ","))
		utils.CheckErr(err)
		for _, g := range groupDetail {
			utils.PrintGroupDetail(g)
		}
	} else if o.deleteGroups{
		groups := strings.Split(o.groups, ",")
		for _, g := range groups {
			utils.CheckErr(admin.DeleteConsumerGroup(g))
			log.Info("Delete consumer group success", zap.String("group", g))
		}
	}else if o.listConsumerOffsets{
		groups := strings.Split(o.groups, ",")
		topics := strings.Split(o.topics, ",")
		_partitions := strings.Split(o.partitions, ",")
		var partitions []int32
		for _, p := range _partitions {
			partition, _ := strconv.Atoi(p)
			partitions = append(partitions, int32(partition))
		}
		topicPartitions := map[string][]int32{}
		for _, t := range topics {
			topicPartitions[t] = partitions
		}
		for _, g := range groups {
			res, err := admin.ListConsumerGroupOffsets(g, topicPartitions)
			utils.CheckErr(err)
			utils.PrintConsumerGroupOffsets(g, res)
		}
	}else if o.describeCluster{
		brokers, controllerID, err := admin.DescribeCluster()
		utils.CheckErr(err)
		utils.PrintCluster(controllerID, brokers)
	}else if o.describeLogDirs{
		_brokers := strings.Split(o.brokers, ",")
		var brokers []int32
		for _, b := range _brokers {
			broker, err := strconv.Atoi(b)
			utils.CheckErr(err)
			brokers = append(brokers, int32(broker))
		}

		res, err := admin.DescribeLogDirs(brokers)
		utils.CheckErr(err)
		utils.PrintLogDirs(res)
	}else {
		cmd.Help()
	}
}

func NewCmdAdmin() *cobra.Command {
	o := newAdminOptions()

	cmd := &cobra.Command{
		Use:     "admin",
		Short:   "Kafka admin operations",
		Long:    "Admin operations. delete records, list consumer groups, describe groups, list group offsets, delete groups, describe cluster, describe log dirs and etc",
		Example: adminExample,
		Run:     o.run,
	}
	cmd.Flags().StringVarP(&o.bootstrapServers, "bootstrap-server", "b", "localhost:9092", "The Kafka server to connect to.more than one should be separated by commas")
	cmd.Flags().BoolVar(&o.deleteRecords, "delete-records", o.deleteRecords, "Delete record, when specified, topics, partitions and offset should also specified")
	cmd.Flags().BoolVar(&o.listConsumerGroups, "list-consumer-groups", o.listConsumerGroups, "List all consumer groups")
	cmd.Flags().BoolVar(&o.describeGroups, "describe-groups", o.describeGroups, "Describe a certain consumer group,when specified, groups should also specified")
	cmd.Flags().BoolVar(&o.deleteGroups, "delete-groups", o.deleteGroups, "Delete consumer groups, when specified, groups should also specified")
	cmd.Flags().BoolVar(&o.listConsumerOffsets, "list-consumer-offsets", o.listConsumerOffsets, "List consumer offsets, when specified, groups, topics, partitions should also specified, and will use cartesian product of topics and partitions")
	cmd.Flags().BoolVar(&o.describeCluster, "describe-cluster", o.describeCluster, "Get information about the nodes in the cluster")
	cmd.Flags().BoolVar(&o.describeLogDirs, "describe-log-dirs", o.describeLogDirs, "Get information about all log directories on the given set of brokers, when specified, brokers should also specified")

	cmd.Flags().StringVar(&o.topics, "topics", o.topics, "The topics commands will act on, separate by commas.")
	cmd.Flags().StringVar(&o.partitions, "partitions", o.partitions, "The partitions commands will act on, separate by commas.")
	cmd.Flags().Int64Var(&o.offset, "offset", o.offset, "The offset commands will act on.")
	cmd.Flags().StringVar(&o.groups, "groups", o.groups, "The consumer groups commands will act on.")
	cmd.Flags().StringVar(&o.brokers, "brokers", o.brokers, "The brokers commands will act on.")
	return cmd
}
