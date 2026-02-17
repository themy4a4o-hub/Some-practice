package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pract/service"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "NiceApp",
		Usage: "ToDoSomeWork",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "path to input file",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "path to output file",
				Value:   "C:/Users/BTC/Desktop/Go/pract/output.txt",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Usage: "debug, info, warn, error",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			logerctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			defer stop()
			input := c.String("input")
			output := c.String("output")
			if output == "" {
				output = "output.txt"
			}
			logLevel := c.String("log-level")
			setupLogger(logLevel, false)
			slog.Info("app started")
			prod := service.Fileproduser(input)
			pres := service.FilePresenter(output)
			srv := service.NewService(prod, pres)
			err := srv.Run(logerctx)
			if err != nil {
				slog.Error("service failed", "err", err)
			} else {
				slog.Info("service finished successfully")
			}
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("app error", "err", err)
	}
}

func setupLogger(levelStr string, jsonOutput bool) {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	var handler slog.Handler
	if jsonOutput {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}
	slog.SetDefault(slog.New(handler))
	slog.Debug("Логгер настроен",
		"level", levelStr,
		"json output", jsonOutput)
}
