package config

import (
	"log"
	"testing"
)

func TestParseConfig(t *testing.T) {
	config := ParseConfig()

	log.Printf("Config %v", config)
}
