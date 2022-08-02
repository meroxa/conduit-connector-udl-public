package udl

import (
	"context"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Destination struct {
	sdk.UnimplementedDestination
	config Config
	client Client
}

func NewDestination() sdk.Destination {
	return &Destination{}
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	parsedCfg, err := Parse(cfg)
	if err != nil {
		return err
	}
	d.config = parsedCfg
	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	c, err := NewClient(d.config.BaseURL)
	if err != nil {
		return err
	}
	d.client = c
	return nil
}

func (d *Destination) WriteAsync(ctx context.Context, record sdk.Record, ackFunc sdk.AckFunc) error {
	return nil
}
