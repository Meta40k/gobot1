package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	APIID         int
	APIHash       string
	Phone         string
	CloudPassword string
}

func LoadConfig() (*Config, error) {

	apiid, _ := strconv.Atoi(os.Getenv("apiID"))
	if apiid <= 0 {
		return nil, fmt.Errorf("apiID must be a positive integer")
	}
	fmt.Println("Config API ID:", apiid)

	apihash := os.Getenv("apiHash")
	if apihash == "" {
		return nil, fmt.Errorf("api hash is required")
	}
	fmt.Println("Config APIHash:", apihash)

	phone := os.Getenv("phone")
	if phone == "" {
		return nil, fmt.Errorf("phone not set")
	}
	fmt.Println("Config phone", phone)

	cloudPassword := os.Getenv("cloudPassword")
	if cloudPassword == "" {
		return nil, fmt.Errorf("пустой облачный пароль")
	}
	fmt.Println("Config cloudPassword", cloudPassword)

	return &Config{
		APIID:         apiid,
		APIHash:       apihash,
		Phone:         phone,
		CloudPassword: cloudPassword,
	}, nil
}
