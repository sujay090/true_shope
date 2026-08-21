package main

import (
	"fmt"
	"log"
	"os"
	"true_shop/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down")
	}

	m, err := migrate.New("file://migrations", cfg.DatabaseUrl)

	if err != nil {
		log.Fatal("migration.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Up successful")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Down successful")
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
	fmt.Println("migration running")

}
