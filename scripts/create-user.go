package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var username, password string
	
	flag.StringVar(&username, "username", "", "Username for the new user")
	flag.StringVar(&password, "password", "", "Password for the new user")
	flag.Parse()

	if username == "" || password == "" {
		fmt.Println("Usage: go run create-user.go --username=<username> --password=<password>")
		os.Exit(1)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password Hash: %s\n", string(hash))
	fmt.Println("TODO: Add database connection and save user")
}
