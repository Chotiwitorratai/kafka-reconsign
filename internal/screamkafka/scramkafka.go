package scramkafka

import (
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"github.com/xdg/scram"
)

var (
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func NewConsumerClient(isAuth bool) (sarama.ConsumerGroup, error) {
	if len(viper.GetString("kafka.group")) == 0 {
		fmt.Sprintln("no Kafka consumer group defined, please set the -group flag")
	}

	version := viper.GetString("kafka.version")
	consumeTopic := viper.GetStringSlice("kafka.topic.consumeTopic")
	brokers := viper.GetStringSlice("kafka.brokers")
	group := viper.GetString("kafka.group")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	kafkaVersion, _ := sarama.ParseKafkaVersion(version)
	config.Version = kafkaVersion

	// log.Println("[CONFIG] [KAFKA] topic:", consumeTopic, "brokers:", brokers, "group:", group)

	log.Printf("[CONFIG] [KAFKA] topic:%s brokers:%s group:%s", consumeTopic, brokers, group)

	if isAuth {
		kafkaUsername := viper.GetString("secrets.kafkaUsername")
		kafkaPassword := viper.GetString("secrets.kafkaPassword")

		config.Net.SASL.Enable = true
		config.Net.SASL.User = kafkaUsername
		config.Net.SASL.Password = kafkaPassword
		config.Net.SASL.Handshake = true

		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }

		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         0,
		}
	}

	return sarama.NewConsumerGroup(viper.GetStringSlice("kafka.brokers"), viper.GetString("kafka.group"), config)
}
