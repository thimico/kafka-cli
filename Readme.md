# Kafka-cli

A command line tool for the apache kafka, include the topic, consumer, producer, admin operations

## Features
- **Topics**
    - list
    - describe topic partitions,replicas
    - create
    - delete
    - add partitions
- **Producer**
    - produce to specify partition
    - produce by specify key
    - produce with headers
    
- **Consumer**
    - consume from specified partition and offset
    
- **ConsumerGroup**
    - consume by group

- **Admin**
    - delete consumer groups
    - delete records
    - describe cluster
    - describe consumer groups
    - describe log dirs
    - list consumer groups
    - list consumer offset

## Installation

    git clone https://github.com/thimico/kafka-cli.git
    cd kafka-cli
    go build cmd/kafka-cli.go

## Usage
**Overview**:

    ./kafka-cli
    
    A command line tools for apache kafka, include topic,consumer,producer, admin's operations
    
    Usage:
      kafka-cli [flags]
      kafka-cli [command]
    
    Available Commands:
      admin       
      consume     Consume kafka message with given topic and group_id
      help        Help about any command
      topic       Kafka topic operations
    
    Flags:
      -h, --help   help for kafka-cli
    
    Use "kafka-cli [command] --help" for more information about a command.
    
**Topic**:
    
    ./kafka-cli.go topic -h   
                                       
    Topic operations, include topic create、list、delete、detail, topic partition create
    
    Usage:
      kafka-cli topic [flags]
    
    Examples:
    
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
    
    
    Flags:
          --add-partition string      The Topic which need to create partition, partition num must higher than which already exists
      -b, --bootstrap-server string   The Kafka server to connect to.more than one should be separated by commas (default "localhost:9092")
      -c, --create string             Create a new topic.
      -d, --delete string             Delete a topic.
          --describe string           List details for the given topics.more than one should be separated by commas
      -h, --help                      help for topic
      -l, --list                      List all available topics.
          --partition-num int32       The specified partition when create topic or add partition (default 1)
          --replica-num int16         The specified replica when create topic (default 1)

**Producer**:
    
    ./kafka-cli.go producer --help
    
    A kafka synchronous producer, with pretty much config options. but it's not asynchronous, which means it will wait for result before return
    
    Usage:
      kafka-cli producer [flags]
    
    Examples:
    
    # Produce a message
        ./kafka-cli producer --bootstrap-servers=localhost:9092 --key=13 --partitioner=random --topic=singed --value='test value'
        result:
            {"level":"info","ts":1612429377.79058,"caller":"log/log.go:16","msg":"Send message success","partition":2,"offset":0}
    
    
    Flags:
      -b, --bootstrap-servers string   The Kafka server to connect to.more than one should be separated by commas (default "localhost:9092")
          --headers string             The headers of the message. Example: -headers=foo:bar,bar:foo
      -h, --help                       help for producer
          --key string                 the key of message
          --partition int32            The partition which message produce to, if provided, it will use manual partitioner (default -1)
          --partitioner string         The partitioning scheme to use. Can be hash, manual, or random (default "hash")
          --topic string               REQUIRED: The topic id to produce messages to.
          --value string               REQUIRED: The message content which is going to be produced
          
**Consumer**
    
    ./kafka-cli consumer -h
    
    Consume kafka message with given topic and partition, if you want to consume all partition, please refer to use consumerg
    
    Usage:
      kafka-cli consumer [flags]
    
    Examples:
    
    # Consume msgs which start with given topic, partition and offset.
        ./kafka-cli consumer --bootstrapServers=localhost:9092 --topic=singed --offset=1 --partition=1
    
    
    Flags:
      -b, --bootstrap-servers string   The Kafka server to connect to.more than one should be separated by commas (default "localhost:9092")
      -h, --help                       help for consumer
          --offset int                 Which offset to consume start with, -2 means oldest, -1 means newest (default -1)
          --partition int32            The partition to consume (default 0)
          --topic string               REQUIRED: The topics to consume,more than one should be separated by commas
          
**Consumer Group**
    
    ./kafka-cli consumerg -h
    
    Consume kafka message with given topics and group_id, and message will be auto committed
    
    Usage:
      kafka-cli consumerg [flags]
    
    Examples:
    
    # Consume kafka messages of set of topic which certain group id, and will print those message in terminal
            kafka-cli consume -t test,singed -c default -b localhost:9092
                    
    
    Flags:
      -b, --bootstrap-servers string   The Kafka server to connect to.more than one should be separated by commas (default "localhost:9092")
          --group-id string            The consumer group ID (default "kafka-cli")
      -h, --help                       help for consumerg
          --topics string              The topics to consume,more than one should be separated by commas
          
**Admin**
    
    ./kafka-cli admin -h
    
    Admin operations. delete records, list consumer groups, describe groups, list group offsets, delete groups, describe cluster, describe log dirs and etc
    
    Usage:
      kafka-cli admin [flags]
    
    Examples:
    
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
    
    
    Flags:
      -b, --bootstrap-server string   The Kafka server to connect to.more than one should be separated by commas (default "localhost:9092")
          --brokers string            The brokers commands will act on.
          --delete-groups             Delete consumer groups, when specified, groups should also specified
          --delete-records            Delete record, when specified, topics, partitions and offset should also specified
          --describe-cluster          Get information about the nodes in the cluster
          --describe-groups           Describe a certain consumer group,when specified, groups should also specified
          --describe-log-dirs         Get information about all log directories on the given set of brokers, when specified, brokers should also specified
          --groups string             The consumer groups commands will act on.
      -h, --help                      help for admin
          --list-consumer-groups      List all consumer groups
          --list-consumer-offsets     List consumer offsets, when specified, groups, topics, partitions should also specified, and will use cartesian product of topics and partitions
          --offset int                The offset commands will act on.
          --partitions string         The partitions commands will act on, separate by commas.
          --topics string             The topics commands will act on, separate by commas.

Please use `./kafka-cli -h` or `./kafka-cli [command] -h` for more detail.

## Compatibility
- **Tested on**
    - apache kafka 2.13.0
    - golang 1.15.7