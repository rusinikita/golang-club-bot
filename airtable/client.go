package airtable

import (
	"context"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
)

type Airtable interface {
	List(ctx context.Context, list interface{}, options ...Options) error
	Create(ctx context.Context, entity any) error
	Patch(ctx context.Context, entity any) error
	Delete(ctx context.Context, entity any) error
}

type config struct {
	Debug  bool   `env:"DEBUG"`
	APIKey string `env:"API_KEY,notEmpty"`
	BaseID string `env:"BASE_ID,notEmpty"`
}

type airtable struct {
	client *resty.Client
}

func New() Airtable {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	client := resty.New()

	client.SetDebug(cfg.Debug)
	client.SetBaseURL("https://api.airtable.com/v0/" + cfg.BaseID)
	client.SetAuthScheme("Bearer")
	client.SetAuthToken(cfg.APIKey)

	return &airtable{client: client}
}

type records struct {
	Records []record `json:"records"`
}

type record struct {
	ID     string         `json:"id,omitempty"`
	Fields map[string]any `json:"fields"`
}
