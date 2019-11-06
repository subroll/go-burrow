package burrow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// NewClient creates new client pointer returned as Client interface
func NewClient(address string, opts ...Option) Client {
	c := new(config)
	defaults(c)
	for _, fn := range opts {
		fn(c)
	}

	return &client{
		config:  c,
		address: strings.TrimRight(address, "/"),
	}
}

type client struct {
	*config
	address string
}

//nolint:errcheck
func (c *client) sendGet(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *client) Ping(ctx context.Context) error {
	b, err := c.sendGet(ctx, c.address+"/burrow/admin")
	if err != nil {
		return err
	}

	if string(b) != "GOOD" {
		return errors.New("burrow is not healthy")
	}

	return nil
}

func (c *client) Clusters(ctx context.Context) ([]string, error) {
	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s", c.address, basePath))
	if err != nil {
		return nil, err
	}

	var resp struct {
		baseResponse
		Clusters []string `json:"clusters"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Clusters, nil
}

func (c *client) Consumers(ctx context.Context, cluster string) ([]string, error) {
	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s/consumer", c.address, basePath, cluster))
	if err != nil {
		return nil, err
	}

	var resp struct {
		baseResponse
		Consumers []string `json:"consumers"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Consumers, nil
}

func (c *client) Topics(ctx context.Context, cluster string) ([]string, error) {
	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s/topic", c.address, basePath, cluster))
	if err != nil {
		return nil, err
	}

	var resp struct {
		baseResponse
		Topics []string `json:"topics"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Topics, nil
}

func (c *client) Cluster(ctx context.Context, name string) (Cluster, error) {
	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s", c.address, basePath, name))
	if err != nil {
		return Cluster{}, err
	}

	var resp struct {
		baseResponse
		Module Cluster `json:"module"`
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return Cluster{}, nil
	}

	return resp.Module, nil
}
