package models

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Task struct {
	Id          string    `json:"id"`
	ParentID    string    `json:"parent_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageQueueConfig struct {
	// MessageQueueConfig represents the message queue-related configuration
	Type     string          `yaml:"type"`
	RabbitMQ *RabbitMQConfig `yaml:"rabbitmq,omitempty"`
	Kafka    *KafkaConfig    `yaml:"kafka,omitempty"`
	Logger   Logger
	Registry *prometheus.Registry
}

// RabbitMQConfig represents the RabbitMQ-specific configuration
type RabbitMQConfig struct {
	URL        string `yaml:"url"`
	Exchange   string `yaml:"exchange"`
	Queue      string `yaml:"queue"`
	RoutingKey string `yaml:"routing_key"`
}

// KafkaConfig represents the Kafka-specific configuration
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}

// LoggingConfig represents the logging-related configuration
type LoggingConfig struct {
	Type        string `yaml:"type"`
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"logLevel"`
}

type Config struct {
	MessageQueue     *MessageQueueConfig `yaml:"messageQueue"`
	Logging          *LoggingConfig      `yaml:"logging"`
	PrometheusConfig *PrometheusConfig   `yaml:"prometheus"`
}

type PrometheusConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
