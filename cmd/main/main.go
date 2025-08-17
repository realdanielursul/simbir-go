package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/realdanielursul/simbir-go/config"
	"github.com/realdanielursul/simbir-go/internal/repository"
	"github.com/realdanielursul/simbir-go/pkg/postgres"
)

// set IDs to one standart

// question: pass time values from service layer
// or delegate it to database default values

// type or *type for null rows in db

// deep think what can admin change in rents

// implement validations

// philosophy: admin methods do not require validation

// add set availability method to transport

func main() {
	cfg, err := config.NewConfig("./config/local.yaml")
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("error creating postgres database: %s", err.Error())
	}

	repositories := repository.NewRepositories(db)
	ctx := context.Background()

	fmt.Println(repositories.Account.List(ctx, 100, 0))
	// fmt.Println(repositories.Payment.UpdateBalance(ctx, 1, -50000))

	// deps := service.ServicesDependencies{
	// 	Repos:    repositories,
	// 	Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
	// 	SignKey:  cfg.JWT.SignKey,
	// 	TokenTTL: cfg.JWT.TokenTTL,
	// }

	// services := service.NewServices(deps)
}
