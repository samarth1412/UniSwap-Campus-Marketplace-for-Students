package config

import (
	"context"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestOpenDBContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := OpenDB(ctx, "host=localhost port=5432 user=postgres password=postgres dbname=uniswap sslmode=disable")
	if err == nil {
		t.Fatalf("expected error for canceled context")
	}
	if !strings.Contains(err.Error(), "ping database") {
		t.Fatalf("expected ping database wrapped error, got %v", err)
	}
}

func TestOpenDBInvalidDSN(t *testing.T) {
	_, err := OpenDB(context.Background(), "%%%")
	if err == nil {
		t.Fatalf("expected invalid dsn error")
	}
}
