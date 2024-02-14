package main

//go:generate sqlc generate -f ../sqlc.yml

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"

	"github.com/dfalgout/dream/config"
	"github.com/dfalgout/dream/custom"
	"github.com/dfalgout/dream/db"
	"github.com/dfalgout/dream/handler"
	"github.com/dfalgout/dream/migrate"
	"github.com/dfalgout/dream/service"
	"github.com/dfalgout/dream/types"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	envVars := config.NewConfig()
	conn, err := sql.Open("sqlite3", "test.db")
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

	e := echo.New()
	e.JSONSerializer = custom.NewJSONSerializer()
	e.Validator = custom.NewValidator()
	e.HideBanner = true
	if envVars.Env == "local" {
		e.Debug = true
	} else {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.JwtClaims)
		},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
		Skipper: func(c echo.Context) bool {
			publicRoutes := []string{handler.ACTION_SEND_CODE, handler.ACTION_VERIFY_CODE}
			return slices.Contains(publicRoutes, c.Request().URL.Path)
		},
	}
	e.Use(echojwt.WithConfig(jwtConfig))

	// services
	userService := service.NewUserService(queries, logger)

	// register handlers
	handler.RegisterUserHandlers(e, logger, userService)

	e.Logger.Fatal(e.Start(":" + envVars.Port))
}
