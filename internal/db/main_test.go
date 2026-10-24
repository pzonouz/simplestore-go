package sqlc

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_store_go?sslmode=disable"
)

var queries *Queries

func TestMain(m *testing.M) {
	ctx := context.TODO()
	conn, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	queries = New(conn)
	os.Exit(m.Run())
}
