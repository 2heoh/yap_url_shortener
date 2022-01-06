package config

import (
	"flag"
	"fmt"
)

var (
	serverAddress   = "localhost:8080"
	baseURL         = fmt.Sprintf("http://%s", serverAddress)
	fileStoragePath = "./links.db"
)

func LoadArgs() (*Config, error) {
	address := new(NetAddress)
	flag.Var(address, "a", "Server address host:port")
	bURL := flag.String("b", baseURL, "Base Url")
	path := flag.String("f", fileStoragePath, "file storage path")
	flag.Parse()

	if address.String() == ":0" {
		address.Host = "localhost"
		address.Port = 8080
	}

	return &Config{
		ServerAddress:   address.String(),
		BaseURL:         *bURL,
		FileStoragePath: *path,
	}, nil
}
