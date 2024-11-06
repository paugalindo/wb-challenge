package bootstrap

import "os"

type Config struct {
	ServicePort      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	NatsURL          string
}

func GetConfigFromEnv() Config {
	return Config{
		ServicePort:      getEnvOrDefault("SVC_PORT", "80"),
		DatabaseHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DatabasePort:     getEnvOrDefault("DB_PORT", "5432"),
		DatabaseUser:     getEnvOrDefault("DB_USER", "wbuser"),
		DatabasePassword: getEnvOrDefault("DB_PASS", "wbpass"),
		DatabaseName:     getEnvOrDefault("DB_NAME", "wbdb"),
		NatsURL:          getEnvOrDefault("NATS_URL", "nats://127.0.0.1:4222"),
	}
}

func getEnvOrDefault(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
