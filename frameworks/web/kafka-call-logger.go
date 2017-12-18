package web

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/urfave/negroni"
)

func (r *RestAPI) KafkaRecorder(response http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	start := time.Now()
	// do some stuff before
	next(response, request)
	// do some stuff after - like logging call
	if r.Registry.Configuration.KafkaMetrics && !r.KafkaInitialized {
		// Set up kafka connection
		r.IPAddress = myIPAddress()

		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Retry.Max = 5
		config.Producer.Return.Successes = true

		broker := fmt.Sprintf("%v:%v", r.Registry.Configuration.KafkaHost, r.Registry.Configuration.KafkaPort)
		brokers := []string{broker}
		msg := fmt.Sprintf("Connecting to Kafka Broker at %v", broker)
		r.Registry.Logger.Log("INFO", msg)
		producer, err := sarama.NewSyncProducer(brokers, config)
		if err != nil {
			panic(err)
		}
		r.Producer = producer
		r.KafkaInitialized = true
	}
	// log call when appropriate
	if r.KafkaInitialized {
		res := response.(negroni.ResponseWriter)
		text := fmt.Sprintf("{ \"Application\":\"%v\",\"Host\":\"%v\",\"Status\":%v,\"Method\":\"%v\",\"Endpoint\":\"%v\",\"Timestamp\":\"%v\",\"Elapsed\":\"%v\" }",
			r.Registry.Configuration.Application,
			r.IPAddress,
			res.Status(),
			request.Method,
			request.URL.Path,
			start.Format(time.RFC3339),
			time.Since(start))
		msg := &sarama.ProducerMessage{
			Topic: r.Registry.Configuration.KafkaMetricsTopic,
			Value: sarama.StringEncoder(text),
		}
		_, _, err := r.Producer.SendMessage(msg)
		if err != nil {
			fmt.Println(err)
		}
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
