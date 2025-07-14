package config

import (
	"context"

	"github.com/Drivello/Twittah/services/auth/ent"
)

func InitEntClient(postgresDSN string) (*ent.Client, error) {
	client, err := ent.Open("postgres", postgresDSN)
	if err != nil {
		return nil, err
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
