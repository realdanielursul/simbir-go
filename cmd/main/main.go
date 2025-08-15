package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/realdanielursul/simbir-go/config"
	"github.com/realdanielursul/simbir-go/internal/repository"
	"github.com/realdanielursul/simbir-go/internal/service"
	"github.com/realdanielursul/simbir-go/pkg/hasher"
	"github.com/realdanielursul/simbir-go/pkg/postgres"
)

// set IDs to one standart

// question: pass time values from service layer
// or delegate it to database default values

// type or *type for null rows in db

// deep think what can admin change in rents

// implement validations

// philosophy: admin methods do not require validation

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

	deps := service.ServicesDependencies{
		Repos:    repositories,
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}

	services := service.NewServices(deps)
	ctx := context.Background()

	fmt.Println(services.AdminAccount.GetAccount(ctx, 1))
}
