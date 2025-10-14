package util

import (
	"log"
	"os"
)

func GetEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("%s is empty or not defined. Defaulting to: %s", key, def)
		v = def
	}
	return v
}
