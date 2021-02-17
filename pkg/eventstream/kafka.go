package eventstream

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func GetKafkaConsumer(groupID string, topics ...string) (*kafka.Consumer, error) {
	config, err := GetConsumerConfig(groupID)
	if err != nil {
		return nil, err
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return consumer, nil
}

func GetConsumerConfig(groupID string) (kafka.ConfigMap, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	config["group.id"] = groupID

	autoOffsetReset, err := mustGet("AUTO_OFFSET_RESET")
	if err != nil {
		return nil, err
	}
	config["auto.offset.reset"] = autoOffsetReset

	return config, nil
}

func GetConfig() (kafka.ConfigMap, error) {
	config := kafka.ConfigMap{}

	bootstrapServers, err := mustGet("BOOTSTRAP")
	if err != nil {
		return nil, err
	}
	config["bootstrap.servers"] = bootstrapServers

	saslMechanishms, err := mustGet("SASL_MECHANISM")
	if err != nil {
		return nil, err
	}
	config["sasl.mechanisms"] = saslMechanishms

	securityProtocol, err := mustGet("SECURITY_PROTOCOL")
	if err != nil {
		return nil, err
	}
	config["security.protocol"] = securityProtocol

	sslCALocation, ok := os.LookupEnv("SSL_CA_LOCATION")
	if ok {
		config["ssl.ca.location"] = sslCALocation
	}

	saslUsername, err := mustGet("SASL_USERNAME")
	if err != nil {
		return nil, err
	}
	config["sasl.username"] = saslUsername

	saslPassword, err := mustGet("SASL_PASSWORD")
	if err != nil {
		return nil, err
	}
	config["sasl.password"] = saslPassword

	return config, nil
}

func GetProducerConfig() (kafka.ConfigMap, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func mustGet(key string) (string, error) {
	data, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.New(key + " required")
	}
	return data, nil
}
