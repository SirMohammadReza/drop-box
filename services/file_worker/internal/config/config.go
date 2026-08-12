package config

type Config struct {
	MinioURL    string `env:"MINIO_URL,required"`
	MinioPort   int    `env:"MINIO_PORT,required"`
	MinioID     string `env:"MINIO_ID,required"`
	MinioSecret string `env:"MINIO_SECRET,required"`

	NatsURL  string `env:"NATS_URL,required"`
	NatsPort int    `env:"NATS_PORT,required"`
}
