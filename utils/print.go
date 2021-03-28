package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

func PrintTopic(topic string, detail sarama.TopicDetail) {
	s, _ := json.MarshalIndent(detail, "", "")
	printSeparator()
	fmt.Printf("TOPIC:%s\nDETAIL:%s\n", topic, string(s))
}

func PrintTopicMeta(meta *sarama.TopicMetadata) {
	printSeparator()
	s, _ := json.MarshalIndent(meta, "", " ")
	fmt.Printf("TOPIC:%s\nDETAIL:%10s\n", meta.Name, string(s))
}

func PrintConsumerMessage(msg *sarama.ConsumerMessage) {
	printSeparator()
	headers:= map[string]string{}
	for _,h:= range msg.Headers {
		headers[string(h.Key)] = string(h.Value)
	}
	fmt.Printf("Headers       :%s\n", headers)
	fmt.Printf("Timestamp     :%s\n", msg.Timestamp)
	fmt.Printf("BlockTimestamp:%s\n", msg.BlockTimestamp)
	fmt.Printf("Key           :%s\n", msg.Key)
	fmt.Printf("Value         :%s\n", msg.Value)
	fmt.Printf("Topic         :%s\n", msg.Topic)
	fmt.Printf("Partition     :%d\n", msg.Partition)
	fmt.Printf("Offset        :%d\n", msg.Offset)
}

func PrintConsumerGroups(groups map[string]string) {
	for k,v := range groups {
		printSeparator()
		fmt.Printf("Name        :%s\n", k)
		fmt.Printf("ProtocolType:%s\n", v)
	}
}

func PrintGroupDetail(g *sarama.GroupDescription) {
	printSeparator()
	fmt.Printf("GroupID       :%s\n", g.GroupId)
	fmt.Printf("State         :%s\n", g.State)
	fmt.Printf("ProtocolType  :%s\n", g.ProtocolType)
	fmt.Printf("Protocol      :%s\n", g.Protocol)
	fmt.Printf("Members       :%s\n", "")
	for k,v := range g.Members {
		fmt.Println("    ------------")
		fmt.Printf("    SessionID  :%s\n", k)
		fmt.Printf("    ClientID   :%s\n", v.ClientId)
		fmt.Printf("    ClientHout :%s\n", v.ClientHost)
		fmt.Printf("    Metadata   :%s\n", string(v.MemberMetadata))
		fmt.Printf("    Assignment :%s\n", string(v.MemberAssignment))
	}
}

func PrintConsumerGroupOffsets(group string, res *sarama.OffsetFetchResponse) {
	printSeparator()
	fmt.Printf("GroupID        :%s\n", group)
	fmt.Printf("Version        :%d\n", res.Version)
	fmt.Printf("ThrottleTimeMs :%d\n", res.ThrottleTimeMs)
	if res.Err != sarama.ErrNoError {
		fmt.Printf("Err            :%s\n", res.Err)
	}
	for k, v := range res.Blocks {
		for k1,v1 := range v {
			fmt.Println("    ------------")
			fmt.Printf("    Topic       :%s\n", k)
			fmt.Printf("    Partition   :%d\n", k1)
			fmt.Printf("    LeaderEpoch :%d\n", v1.LeaderEpoch)
			fmt.Printf("    Metadata    :%s\n", v1.Metadata)
			fmt.Printf("    Offset      :%d\n", v1.Offset)
		}
	}
}

func PrintCluster(controllerID int32, brokers []*sarama.Broker) {
	printSeparator()
	fmt.Println("ControllerID:", controllerID)
	for _, b := range brokers {
		fmt.Printf("%-10s%-30s\n", "BrokerID", "BrokerAddr")
		fmt.Printf("%-10d%-30s\n",b.ID(), b.Addr())
	}

}

func PrintLogDirs(info map[int32][]sarama.DescribeLogDirsResponseDirMetadata) {
	for k, v := range info {
		printSeparator()
		fmt.Printf("BrokerID:%d\n", k)
		for _, d := range v {
			fmt.Printf("Path: %s\n", d.Path)
			fmt.Printf("%-50s%-15s%-20s%-20s%-20s\n", "Topic", "IsTemporary", "OffsetLag", "PartitionID", "Size")
			for _, t := range d.Topics {
				for _, p := range t.Partitions {
					fmt.Printf("%-50s%-15t%-20d%-20d%-20d\n",t.Topic, p.IsTemporary, p.OffsetLag, p.PartitionID, p.Size)
				}
			}
		}
	}
}

func printSeparator() {
	fmt.Println("*****************************************************")
}
