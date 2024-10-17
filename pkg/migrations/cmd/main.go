package main

import (
	"log"
	"os"

	"github.com/wesleymassine/swordhealth/migrations"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing argument: please specify 'migrate-up' or 'migrate-down'")
	}

	command := os.Args[1]
	switch command {
	case "migrate-up":
		migrations.RunMigrationsUp()
	case "migrate-down":
		migrations.RunMigrationsDown()
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
