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
	"github.com/realdanielursul/simbir-go/pkg/logger"
	"github.com/realdanielursul/simbir-go/pkg/postgres"
)

// set IDs to one standart

// deep think what can admin change in rents

// implement validations

// philosophy: admin methods do not require validation

func main() {
	logger.SetLogrus()

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

	go services.Payment.BillingWorker(ctx)

	fmt.Println(services.AdminAccount.CreateAccount(ctx, &service.AdminAccountInput{
		Username: "danixx",
		Password: "c094yt178",
		IsAdmin:  true,
		Balance:  25000,
	}))

	fmt.Println(services.Transport.CreateTransport(ctx, 1, &service.TransportInput{
		CanBeRented:   true,
		TransportType: "Car",
		Model:         "model",
		Color:         "color",
		Identifier:    "identifier",
		Description:   nil,
		Latitude:      15.656,
		Longitude:     47.123,
		MinutePrice:   15,
		DayPrice:      15000,
	}))

	fmt.Println(services.Rent.StartRent(ctx, 1, 1, "Minutes"))

	select {}
}
