//nolint:golint
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

var _ Client = (*client)(nil)

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

func (c *client) Clusters(ctx context.Context) (Clusters, error) {
	var resp struct {
		baseResponse
		Clusters Clusters `json:"clusters"`
	}

	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s", c.address, basePath))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Clusters, nil
}

func (c *client) Consumers(ctx context.Context, cluster string) (Consumers, error) {
	var resp struct {
		baseResponse
		Consumers Consumers `json:"consumers"`
	}

	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s/consumer", c.address, basePath, cluster))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Consumers, nil
}

func (c *client) Topics(ctx context.Context, cluster string) (Topics, error) {
	var resp struct {
		baseResponse
		Topics Topics `json:"topics"`
	}

	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s/topic", c.address, basePath, cluster))
	if err != nil {
		return nil, err
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
	var resp struct {
		baseResponse
		Module Cluster `json:"module"`
	}

	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s", c.address, basePath, name))
	if err != nil {
		return resp.Module, err
	}

	if err := json.Unmarshal(b, &resp); err != nil {
		return resp.Module, err
	}

	if resp.Error {
		return resp.Module, errors.New(resp.Message)
	}

	return resp.Module, nil
}

func (c *client) Consumer(ctx context.Context, cluster, name string) (Consumer, error) {
	var resp struct {
		baseResponse
		Topics Consumer `json:"topics"`
	}

	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s/consumer/%s", c.address, basePath, cluster, name))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Topics, nil
}

func (c *client) Topic(ctx context.Context, cluster, name string) (Topic, error) {
	var resp struct {
		baseResponse
		Offsets Topic `json:"offsets"`
	}

	b, err := c.sendGet(ctx, fmt.Sprintf("%s%s/%s/topic/%s", c.address, basePath, cluster, name))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New(resp.Message)
	}

	return resp.Offsets, nil
}

// nolint: errcheck
func (c *client) DeleteConsumer(ctx context.Context, cluster, consumer string) error {
	path := fmt.Sprintf("%s%s/%s/topic/%s", c.address, basePath, cluster, consumer)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		baseResponse
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if r.Error {
		return errors.New(r.Message)
	}

	return nil
}
