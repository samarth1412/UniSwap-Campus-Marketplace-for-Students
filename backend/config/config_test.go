package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "secret",
		DBName:     "uniswap",
		DBSSLMode:  "disable",
	}

	got := cfg.DatabaseDSN()
	want := "host=localhost port=5432 user=postgres password=secret dbname=uniswap sslmode=disable"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestLoadResourceEnvParsesValues(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, "app.env")
	content := "# comment\nPORT=9090\nDB_PASSWORD=\"abc123\"\nINVALID_LINE\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	values, err := loadResourceEnv(envPath)
	if err != nil {
		t.Fatalf("loadResourceEnv failed: %v", err)
	}

	if values["PORT"] != "9090" {
		t.Fatalf("expected PORT to be parsed")
	}
	if values["DB_PASSWORD"] != "abc123" {
		t.Fatalf("expected DB_PASSWORD to be unquoted")
	}
}

func TestGetConfigValuePrefersEnvOverFile(t *testing.T) {
	t.Setenv("PORT", "1234")
	fileValues := map[string]string{"PORT": "5678"}
	got := getConfigValue(fileValues, "PORT", "8080")
	if got != "1234" {
		t.Fatalf("expected env value, got %q", got)
	}
}

func TestLoadReadsResourceFile(t *testing.T) {
	tmpDir := t.TempDir()
	resourceDir := filepath.Join(tmpDir, "resources")
	if err := os.MkdirAll(resourceDir, 0o755); err != nil {
		t.Fatalf("mkdir resources: %v", err)
	}

	envBody := "PORT=9999\nDB_HOST=127.0.0.1\nDB_USER=test\nDB_PASSWORD=pwd\nDB_NAME=testdb\nDB_SSLMODE=disable\nJWT_SECRET=jwt\n"
	if err := os.WriteFile(filepath.Join(resourceDir, "app.env"), []byte(envBody), 0o644); err != nil {
		t.Fatalf("write app.env: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir tempdir: %v", err)
	}

	t.Setenv("PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_SSLMODE", "")
	t.Setenv("JWT_SECRET", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Port != "9999" || cfg.DBPassword != "pwd" || cfg.JWTSecret != "jwt" {
		t.Fatalf("unexpected config loaded: %+v", cfg)
	}
}
