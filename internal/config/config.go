package config

type Config struct {
	ServerAddress string
	// Add more configuration options here
}

func New() *Config {
	return &Config{
		ServerAddress: ":5000",
		// Initialize other config values
	}
}
