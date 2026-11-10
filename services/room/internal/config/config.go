package config

import "os"

type Config struct {
	GRPCAddr string
}

func Load() Config {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":50051"
	}
	return Config{GRPCAddr: addr}
}
