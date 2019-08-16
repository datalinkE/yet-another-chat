package config

import "os"

type Config struct {
	BoltFile   string
	ListenAddr string
}

func GetFromEnv() *Config {
	boltFile := os.Getenv("BOLT_FILE")
	if boltFile == "" {
		boltFile = "bolt.db"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":9000"
	}

	return &Config{
		BoltFile:   boltFile,
		ListenAddr: listenAddr,
	}
}
