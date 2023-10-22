package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	discord "github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service/internal/user"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

/*
	TODOs:
		META:
			* start a documentation diretcory
			* start implementing testing coverage
			* work more on readme
			* abstract what we can to delphi-go-kit (e.g. logging, tracing, etc.)
			* determine what logging tracing solutions I want to use long term(probably just something within aws honestly)
			* refactor rootlevel protobuf/grpc logic into corresponding
				internal directories
			* refactor main.go to cmd/delphi-x-service/main.go
			* clean up Make targets and keep them up to date
			* abstract build logic execution into submodule delphi build-utils

		MESA:
			* we need to refactor the hell out of the discord bot logic, it is the smelliest dirtiest code this side of the mississippi
			* we also need a cleaner way to integrate user and admin bot logic than the current implementation
*/

func main() {
	ctx := context.Background()
	run(ctx, os.Args)
}

func run(ctx context.Context, _ []string) error {
	var cfg struct {
		ServiceName           string `env:"SERVICE_NAME" envDefault:"delphi-discord-bot-client-service"`
		InfrastructureService struct {
			Host string `env:"INFRASTRUCTURE_SERVICE_HOST" envDefault:"localhost"`
			Port string `env:"INFRASTRUCTURE_SERVICE_PORT" envDefault:"8080"`
		}
		InferentialService struct {
			Host string `env:"INFERENTIAL_SERVICE_HOST" envDefault:"localhost"`
			Port string `env:"INFERENTIAL_SERVICE_PORT" envDefault:"8082"`
		}
		TrainingService struct {
			Host string `env:"INFERENTIAL_SERVICE_HOST" envDefault:"localhost"`
			Port string `env:"INFERENTIAL_SERVICE_PORT" envDefault:"8082"`
		}
		Env string `env:"ENV" envDefault:"local"`
		API struct {
			Address string `env:"API_ADDRESS" envDefault:"http://localhost:80"`
			Port    string `env:"API_PORT" envDefault:"80"`
		}
	}
	if err := env.Parse(&cfg); err != nil {
		return errors.Wrap(err, "parsing configuration")
	}

	discord.InitiateDiscordBotSession(ctx)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}
