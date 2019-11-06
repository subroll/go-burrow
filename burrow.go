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
	Ping(ctx context.Context) error
	Clusters(ctx context.Context) ([]string, error)
	Consumers(ctx context.Context, cluster string) ([]string, error)
	Topics(ctx context.Context, cluster string) ([]string, error)

	Cluster(ctx context.Context, name string) (Cluster, error)
}

type baseResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Request struct {
		URI  string `json:"uri"`
		Host string `json:"host"`
	} `json:"request"`
}

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
