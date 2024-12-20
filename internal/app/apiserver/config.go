package apiserver

type Config struct {
	Address      string `toml:"address"`
	Database_url string `toml:"database_url"`
	Secret_key   string `toml:"secret_key"`
}

func NewConfig() *Config {
	return &Config{
		Address: "localhost:5000",
	}
}
