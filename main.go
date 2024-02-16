package main

//go:generate sqlc generate

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dfalgout/dream/config"
	"github.com/dfalgout/dream/custom"
	"github.com/dfalgout/dream/db"
	"github.com/dfalgout/dream/handler"
	"github.com/dfalgout/dream/migrate"
	"github.com/dfalgout/dream/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "modernc.org/sqlite"
)

//go:embed assets
var assets embed.FS

func main() {
	envVars := config.NewConfig()
	conn, err := sql.Open("sqlite", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	if envVars.AutoMigrations {
		migration := migrate.NewMigration()
		n, err := migration.Up(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Applied %d migrations!\n", n)
	}

	queries := db.New(conn)

	app := echo.New()
	app.JSONSerializer = custom.NewJSONSerializer()
	app.Validator = custom.NewValidator()
	app.HideBanner = true
	if envVars.Env == "local" {
		app.Debug = true
	} else {
		app.Use(middleware.Logger())
	}
	app.Use(middleware.Recover())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      false,
		Root:       "assets",
		Browse:     false,
		Filesystem: http.FS(assets),
	}))
	app.File("/favicon.ico", "assets/favicon.ico")
	app.Use(handler.JwtGuard())
	app.Use(handler.RouteGuard)
	app.RouteNotFound("/*", handler.Fallback)

	// services
	userService := service.NewUserService(queries, logger)

	// register handlers
	handler.RegisterUserHandlers(app, logger, userService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := app.Start(":" + envVars.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
}
