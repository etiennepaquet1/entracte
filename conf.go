package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var ZeffyAPI string

func LoadConfig() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	var ok bool
	ZeffyAPI, ok = os.LookupEnv("ZEFFY_API")
	if !ok || ZeffyAPI == "" {
		return fmt.Errorf("ZEFFY_API is not set")
	}
	return nil
}
