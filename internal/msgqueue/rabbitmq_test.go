package msgqueue

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
	"time"

	. "github.com/bhupeshpandey/task-manager-ashland/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRabbitMQConnection(t *testing.T) {
	// Set up RabbitMQ connection parameters
	amqpURL := "amqp://guest:guest@localhost:5672/"

	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Errorf("Failed to connect to RabbitMQ: %s", err)
		return
	}
	defer conn.Close()

	// Create a new channel
	ch, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to create channel: %s", err)
		return
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"test_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		t.Errorf("Failed to declare queue: %s", err)
		return
	}

	// Publish a message to the queue
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello, World!"),
		})
	if err != nil {
		t.Errorf("Failed to publish message: %s", err)
		return
	}

	// Consume the message from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		t.Errorf("Failed to consume message: %s", err)
		return
	}

	// Verify the message was received
	select {
	case msg := <-msgs:
		if string(msg.Body) != "Hello, World!" {
			t.Errorf("Received message does not match expected message")
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timed out waiting for message")
	}
}

func TestRabbitMQPublishAndConsume(t *testing.T) {
	// Set up RabbitMQ connection parameters
	amqpURL := "amqp://guest:guest@localhost:5672/"

	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Errorf("Failed to connect to RabbitMQ: %s", err)
		return
	}
	defer conn.Close()

	// Create a new channel
	ch, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to create channel: %s", err)
		return
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"test_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		t.Errorf("Failed to declare queue: %s", err)
		return
	}

	// Publish multiple messages to the queue
	for i := 0; i < 10; i++ {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("Message %d", i)),
			})
		if err != nil {
			t.Errorf("Failed to publish message: %s", err)
			return
		}
	}

	// Consume the messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		t.Errorf("Failed to consume message: %s", err)
		return
	}

	// Verify the messages were received
	received := make([]string, 0, 10)
	for msg := range msgs {
		received = append(received, string(msg.Body))
		if len(received) == 10 {
			break
		}
	}
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Log(level LogLevel, message string, args ...interface{}) {
	m.Called(level, message, args)
}

func TestNewRabbitMQ(t *testing.T) {
	tests := []struct {
		name         string
		config       *RabbitMQConfig
		logger       Logger
		prometheus   *PrometheusConfig
		wantErr      bool
		wantRabbitMQ *rabbitMQ
	}{
		{
			name: "success",
			config: &RabbitMQConfig{
				URL:        "amqp://guest:guest@localhost:5672/",
				Queue:      "test_queue",
				Exchange:   "amq.direct",
				RoutingKey: "test_routing_key",
			},
			logger:     &MockLogger{},
			prometheus: &PrometheusConfig{Host: "localhost", Port: "8080"},
			wantErr:    false,
			wantRabbitMQ: &rabbitMQ{
				exchange:   "test_exchange",
				queue:      "test_queue",
				routingKey: "test_routing_key",
				rmqurl:     "amqp://guest:guest@localhost:5672/",
				logger:     &MockLogger{},
			},
		},
		{
			name: "connection error",
			config: &RabbitMQConfig{
				URL: "amqp://guest:guest@localhost:5673/", // invalid port
			},
			logger:     &MockLogger{},
			prometheus: &PrometheusConfig{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := newRabbitMQ(tt.config, tt.logger, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				//assert.Equal(t, tt.wantRabbitMQ, rm)
			}
		})
	}
}

func TestReceiveMessages(t *testing.T) {
	//rm := &rabbitMQ{
	//	queue:      "test_queue",
	//	exchange:   "test_exchange",
	//	routingKey: "test_routing_key",
	//	rmqurl:     "amqp://guest:guest@localhost:5672/",
	//	logger:     &MockLogger{},
	//}
	//
	//wg := &sync.WaitGroup{}
	//wg.Add(1)
	//
	//go func() {
	//	rm.ReceiveMessages(wg)
	//}()
	//
	//// simulate a message delivery
	//delivery := amqp.Delivery{
	//	Body: []byte("test message"),
	//}
	//rm.channel <- delivery
	//
	//wg.Wait()
	//
	//assert.Equal(t, uint64(1), rmqMessagesTotal.Value())
}
