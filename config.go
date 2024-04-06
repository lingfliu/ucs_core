package main

type Config struct {
	nodes []string
}

func NewDefaultConfig() *Config {
	return &Config{
		nodes: []string{"http://localhost:9200"},
	}
}
