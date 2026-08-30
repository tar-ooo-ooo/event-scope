package config

type Config struct {
	Port             string
	Kafka            Kafka
	FrontendEndpoint string
}

type Kafka struct {
	Broker string
	Topics []Topic
}

type Topic struct {
	Name string
}

func Load() Config {
	return Config{
		Port: "8080",
		Kafka: Kafka{
			Broker: "localhost:9092",
			Topics: []Topic{{Name: "events"}, {Name: "events.dlq"}},
		},
		FrontendEndpoint: "http://localhost:5173",
	}
}
