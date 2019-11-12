// Package burrow provides burrow client to interact with burrow HTTP API
package burrow

import (
	"context"
	"net/http"
)

const basePath string = "/v3/kafka"

type config struct {
	httpClient *http.Client
}

func defaults(c *config) {
	c.httpClient = http.DefaultClient
}

// Option represents a function that can be provided as a parameter to NewClient
type Option func(c *config)

// WithHTTPClient sets the HTTP client to be used by the client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) {
		c.httpClient = httpClient
	}
}

// Client is the interface that should be implemented by client struct
type Client interface {
	// Ping is used to check burrow healthiness
	Ping(ctx context.Context) error

	// Clusters will return either list of clusters name or error
	Clusters(ctx context.Context) (Clusters, error)

	// Consumers will return either list of consumers name or error
	Consumers(ctx context.Context, cluster string) (Consumers, error)

	// Topics will return either list of topics name or error
	Topics(ctx context.Context, cluster string) (Topics, error)

	// Cluster will return either cluster detail or error
	Cluster(ctx context.Context, name string) (Cluster, error)

	// Consumer will return either consumer detail or error
	Consumer(ctx context.Context, cluster, name string) (Consumer, error)

	// Topic will return either topic detail or error
	Topic(ctx context.Context, cluster, name string) (Topic, error)

	// DeleteConsumer will delete consumer through burrow and return error if fail
	DeleteConsumer(ctx context.Context, cluster, consumer string) error
}

type baseResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Request struct {
		URI  string `json:"uri"`
		Host string `json:"host"`
	} `json:"request"`
}

// Clusters is the type for cluster list
type Clusters []string

// Consumers is the type for consumer list
type Consumers []string

// Topics is the type for topic list
type Topics []string

// Topic is the type for topic detail
type Topic []int

// Consumer is the type for consumer detail
type Consumer map[string][]ConsumerTopicDetail

// Cluster is the structure for cluster detail
type Cluster struct {
	ClassName     string `json:"class-name"`
	ClientProfile struct {
		ClientID     string `json:"client-id"`
		KafkaVersion string `json:"kafka-version"`
		Name         string `json:"name"`
	} `json:"client-profile"`
	OffsetRefresh int      `json:"offset-refresh"`
	Servers       []string `json:"servers"`
	TopicRefresh  int      `json:"topic-refresh"`
}

// ConsumerTopicDetail is the structure for consumer topic detail
type ConsumerTopicDetail struct {
	CurrentLag int `json:"current-lag"`
	Offsets    []struct {
		Lag       int `json:"lag"`
		Offset    int `json:"offset"`
		Timestamp int `json:"timestamp"`
	} `json:"offsets"`
	Owner string `json:"owner"`
}
