package main

import (
	"context"
	"path"
	"testing"
	"time"

	"github.com/ddbgio/nailivic/db/sqlc"
	"github.com/ddbgio/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	opts := postgres.PostgresOpts{
		Host:     "localhost",
		User:     "testUser",
		Password: "testPassword",
		Name:     "testdb-nailivic",
		Sslmode:  "disable",
	}
	ctx := context.Background()
	db, teardown, err := postgres.NewTestDB(ctx, opts)
	require.NoError(t, err)
	defer teardown()
	pool, err := db.Pool(ctx)
	require.NoError(t, err)
	defer pool.Close()

	err = db.Ping(ctx)
	require.NoError(t, err)

	// test the database
	migrationsDir := path.Join("db", "migrations")
	migrations, err := postgres.Migrations(migrationsDir, "up")
	for _, migration := range migrations {
		t.Logf("applying migration %s", migration.Filename)
		results, err := db.Query(ctx, migration.Content)
		require.NoError(t, err)
		t.Logf("results: %v", results)
	}
	require.NoError(t, err)

	// setup connection pool and queries
	require.NoError(t, err)
	queries := sqlc.New(pool)

	// add a user
	testUsername := "testUser"
	testPassword := "asdfqwerty"
	newUser := sqlc.UserAddParams{
		Username:       testUsername,
		HashedPassword: testPassword,
		CreatedOn: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		LastLogin: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
	queries.UserAdd(ctx, newUser)

	// get the user
	user, err := queries.UserByUsername(ctx, testUsername)
	require.NoError(t, err)
	require.Equal(t, testUsername, user.Username)
	t.Logf("new user: %v", user)

	// check for migration added user
	testAdmin := "admin"
	user, err = queries.UserByUsername(ctx, testAdmin)
	require.NoError(t, err)
	require.Equal(t, testAdmin, user.Username)
	t.Logf("admin user: %v", user)
}
