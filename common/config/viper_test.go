package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type discardLogger struct{}

func (discardLogger) Info(string, ...interface{})    {}
func (discardLogger) Warning(string, ...interface{}) {}
func (discardLogger) Error(string, ...interface{})   {}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestViperConfig_GetConfig(t *testing.T) {
	path := writeConfig(t, `
databases:
  - name: core
    type: mongo
    host: 127.0.0.1
    port: 27017
    database: nutrix
    username: ""
    password: ""
    file_path: ./data/db
    tables:
      sales: sales
auth:
  jwt_secret: "secret"
  jwt_expire_hrs: 24
  enabled: true
zitadel:
  domain: "zitadel"
  port: 8080
  key_path: "./key.json"
  enabled: false
serve_frontend: true
timezone: Africa/Cairo
env: dev
`)

	cfg := ConfigFactory("viper", path, discardLogger{})

	require.Len(t, cfg.Databases, 1)
	db := cfg.Databases[0]
	assert.Equal(t, "mongo", db.Type)
	assert.Equal(t, "127.0.0.1", db.Host)
	assert.Equal(t, 27017, db.Port)
	assert.Equal(t, "nutrix", db.Database)
	assert.Equal(t, map[string]string{"sales": "sales"}, db.Tables)
	assert.Equal(t, "secret", cfg.Auth.JWTSecret)
	assert.Equal(t, 24, cfg.Auth.JWTExpireHrs)
	assert.True(t, cfg.Auth.Enabled)
	assert.Equal(t, "zitadel", cfg.Zitadel.Domain)
	assert.Equal(t, 8080, cfg.Zitadel.Port)
	assert.True(t, cfg.ServeFrontEnd)
	assert.Equal(t, "Africa/Cairo", cfg.TimeZone)
	assert.Equal(t, "dev", cfg.Env)
}

func TestViperConfig_MissingFile(t *testing.T) {
	vc := NewViperConfig(discardLogger{})
	vc.ReadFile(filepath.Join(t.TempDir(), "nope.yaml"))

	_, err := vc.GetConfig()
	assert.Error(t, err)
}
