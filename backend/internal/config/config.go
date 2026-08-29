package config

type Config struct {
	Port             string
	KafkaBrokers     string
	FrontendEndpoint string
}

func Load() Config {
	return Config{
		Port:             "8080",
		KafkaBrokers:     "localhost:9092",
		FrontendEndpoint: "http://localhost:5173",
	}
}
