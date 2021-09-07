package apiserver

//Config ...
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	JwtSignKey  string `toml:"jwtsignkey"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr:   ":8080",
		LogLevel:   "debug",
		JwtSignKey: "",
	}
}
