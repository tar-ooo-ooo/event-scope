package config

type Config struct {
	Port         string
	KafkaBrokers string
}

func Load() Config {
	return Config{
		Port:         "8080",
		KafkaBrokers: "localhost:9092",
	}
}
