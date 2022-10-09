package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	url  = "https://goerli.ethereum.coinbasecloud.net"
	auth = "MTSJPQSCJQU6YWKIGJVG:S3E4XACGWSK26DORYXVUASRSKI2GSZBQQ7CKMKO2"
	POST = "POST"
	GET  = "GET"

	contentTypeJson   = "application/json"
	headerContentType = "Content-Type"
	headerUser        = "user"
)

type CloudClient struct {
	client *http.Client
}

func NewCloudClient(ctx context.Context) *CloudClient {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	client.Get(url)
	return &CloudClient{
		client: client,
	}
}

func (c *CloudClient) setCommonHeaders(req *http.Request) {
	req.Header.Set(headerContentType, contentTypeJson)
	req.Header.Set(headerUser, auth)
}

func (c *CloudClient) GetBalance(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, GET, url, nil)
	if err != nil {
		return err
	}
	c.setCommonHeaders(req)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("Response from Node %s", res)
	return nil
}
