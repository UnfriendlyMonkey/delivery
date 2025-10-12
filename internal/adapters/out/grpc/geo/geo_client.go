// Package geo
package geo

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"
	"delivery/internal/pkg/errs"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.GeoClient = &Client{}

type Client struct {
	conn *grpc.ClientConn
	pbClient geopb.GeoClient
	timeout time.Duration
}

func NewClient(host string) (*Client, error) {
	if host == "" {
		return nil, errs.NewValueIsInvalidError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbClient := geopb.NewGeoClient(conn)

	return &Client{
		conn: conn,
		pbClient: pbClient,
		timeout: 5 * time.Second,
	}, nil
}

func (c *Client) GetGeoLocation(_ context.Context, street string) (kernel.Location, error) {
	request := &geopb.GetGeolocationRequest{
		Street: street,
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.pbClient.GetGeolocation(ctx, request)
	if err != nil {
		return kernel.Location{}, err
	}

	return kernel.NewLocation(uint8(resp.Location.X), uint8(resp.Location.Y))
}

func (c *Client) Close() error {
	return c.conn.Close()
}
