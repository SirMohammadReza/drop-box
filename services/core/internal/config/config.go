package config

type Config struct {
	GRPCURL           string `env:"GRPC_URL,required"`
	GRPCAuthPort      int    `env:"GRPC_AUTH_PORT,required"`
	GRPCIngestionPort int    `env:"GRPC_INGESTION_PORT,required"`

	MongoURL      string `env:"MONGO_URL,required"`
	MongoPort     int    `env:"MONGO_PORT,required"`
	MongoUsername string `env:"MONGO_USERNAME,required"`
	MongoPassword string `env:"MONGO_PASSWORD,required"`
}
