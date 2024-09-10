package main

import (
	"context"
	"testing"

	pg "github.com/ddbgio/postgres"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	opts := pg.PostgresOpts{
		Host:     "localhost",
		User:     "testUser",
		Password: "testPassword",
		Name:     "testdb-nailivic",
		Sslmode:  "disable",
	}
	ctx := context.Background()
	postgres, teardown, err := pg.NewTestDB(ctx, opts)
	require.NoError(t, err)
	defer teardown()
	err = postgres.Ping(ctx)
	require.NoError(t, err)
}
