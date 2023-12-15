package main

import (
	"context"
	"embed"
	"github.kolesa-team.org/backend/go-example/cmd"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
	"os"
	"os/signal"
)

// Вызываем тут вручную, т.к. дефолтный любит писать в логи
func init() { _, _ = maxprocs.Set() }

// подключаем статику в бинарник

//go:embed public/*
var public embed.FS

//go:embed resources/*
var resources embed.FS

func main() {
	// nolint
	ctx, cancel := signal.NotifyContext(
		context.WithValue(
			context.WithValue(
				context.Background(),
				"fs::resources",
				resources,
			),
			"fs::public",
			public,
		),
		os.Interrupt, os.Kill,
	)

	defer cancel()

	// запускаем корневую команду
	if err := cmd.Execute(ctx); err != nil {
		log.Println(err)
	}
}
