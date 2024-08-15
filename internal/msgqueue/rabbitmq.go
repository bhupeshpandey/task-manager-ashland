package msgqueue

import (
	. "github.com/bhupeshpandey/task-manager-ashland/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type rabbitMQ struct {
	exchange   string
	queue      string
	routingKey string
	rmqurl     string
	conn       *amqp.Connection
	channel    *amqp.Channel
	logger     Logger
}

var (
	rmqMessagesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_total_messages",
		Help: "Total number of Rabbit MQ messages received",
	})
)

func newRabbitMQ(config *RabbitMQConfig, logger Logger, registry *prometheus.Registry) (ReceiverMessageQueue, error) {

	if registry != nil {
		registry.MustRegister(rmqMessagesTotal)
	}
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(config.URL)

	if err != nil {
		logger.Log(ErrorLevel, "Failed to connect to RabbitMQ", err)
		return nil, err
	}

	// Create a new channel
	ch, err := conn.Channel()

	if err != nil {
		logger.Log(ErrorLevel, "Failed to open a channel", err)
		return nil, err
	}

	// Declare a queue
	_, err = ch.QueueDeclare(
		config.Queue, // Queue name
		false,        // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)

	if err != nil {
		logger.Log(ErrorLevel, "Failed to declare a queue", err)
		return nil, err
	}

	// Bind the queue to the exchange with a routing key
	routingKey := config.RoutingKey
	err = ch.QueueBind(
		config.Queue,    // Queue name
		routingKey,      // Routing key
		config.Exchange, // Exchange name
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		logger.Log(ErrorLevel, "Failed to bind the queue to the exchange", err)
		return nil, err
	}
	rmq := &rabbitMQ{queue: config.Queue, exchange: config.Exchange, routingKey: config.RoutingKey, rmqurl: config.URL, channel: ch, conn: conn, logger: logger}
	return rmq, nil
}

func (r *rabbitMQ) ReceiveMessages(wg *sync.WaitGroup) error {
	defer wg.Done()
	// Set QoS
	err := r.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		r.logger.Log(ErrorLevel, "Failed to set QoS", err)
		return err
	}

	deliveryChannel, err := r.channel.Consume(r.queue, "", true, true, true, false, nil)
	if err != nil {
		r.logger.Log(ErrorLevel, "Failed to create delivery channel to receive messages", err)
		return err
	}

	for data := range deliveryChannel {
		rmqMessagesTotal.Inc()
		body := data.Body
		strBody := string(body)
		r.logger.Log(InfoLevel, "Received Data ", strBody)
	}

	return nil
}
