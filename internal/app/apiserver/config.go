package apiserver

type Config struct {
	Address string `toml:"address"`
}

func NewConfig() *Config {
	return &Config{
		Address: "localhost:5000",
	}
}
