// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

const (
	DEFAULT_PERIOD             time.Duration = 10 * time.Second
	DEFAULT_HOST               string        = "localhost"
	DEFAULT_PORT               int           = 6379
	DEFAULT_NETWORK            string        = "tcp"
	DEFAULT_MAX_CONN           int           = 10
	DEFAULT_DBID               int           = 0
	DEFAULT_LUA_SCRIPT         string        = "./smshubbeat.lua"
	DEFAULT_AUTH_REQUIRED      bool          = false
	DEFAULT_AUTH_REQUIRED_PASS string        = ""
)

type RedisConfig struct {
	Host          *string
	Port          *int
	Dbid          *int
	Luascript     *string
	Network       *string
	Maxconn       *int
	Auth          struct {
		Required     *bool   `yaml:"required"`
		Requiredpass *string `yaml:"requiredpass"`
	}
}

type BeatConfig struct {
	Period *int64
}

type ConfigSettings struct {
	RedisSettings RedisConfig
	BeatSettings  BeatConfig
}
