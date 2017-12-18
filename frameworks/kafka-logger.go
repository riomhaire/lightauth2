package frameworks

import (
	"fmt"
	"net"
	"os"

	"github.com/Shopify/sarama"
)

type KafaLogger struct {
	Broker    string
	Topic     string
	Producer  sarama.SyncProducer
	IPAddress string
}

func NewKafkaLogger(host string, port int, topic string) *KafaLogger {
	logger := KafaLogger{}
	logger.IPAddress = myIPAddress()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	logger.Broker = fmt.Sprintf("%v:%v", host, port)
	logger.Topic = topic
	brokers := []string{logger.Broker}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}
	logger.Producer = producer
	return &logger
}

func (d KafaLogger) Log(level, message string) {
	text := fmt.Sprintf("[%s][%s] %s\n", d.IPAddress, level, message)
	msg := &sarama.ProducerMessage{
		Topic: d.Topic,
		Value: sarama.StringEncoder(text),
	}
	_, _, err := d.Producer.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
	}

}

func myIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "UnknownHost"

}
