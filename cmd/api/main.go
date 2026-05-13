package main

import (
	"fmt"

	"github.com/Adityoexs/ELA-BE/internal/auth"
	"github.com/Adityoexs/ELA-BE/internal/config"
	"github.com/Adityoexs/ELA-BE/internal/database"
	"github.com/Adityoexs/ELA-BE/internal/employee"
	transporthttp "github.com/Adityoexs/ELA-BE/internal/http"
	"github.com/Adityoexs/ELA-BE/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	log := logger.New(cfg.App.Env, cfg.App.LogLevel)

	db, err := database.New(cfg.Database, log)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}

	if err := db.AutoMigrate(&employee.Employee{}); err != nil {
		log.WithError(err).Fatal("failed to run auto migration")
	}

	repo := employee.NewGormRepository(db)
	svc := employee.NewService(repo, log.WithField("component", "employee_service"))
	endpoints := employee.NewEndpoints(svc)
	handler := employee.NewHandler(endpoints, log.WithField("component", "employee_handler"))

	authSvc := auth.NewService(cfg.JWT)
	authHandler := auth.NewHandler(authSvc, log.WithField("component", "auth_handler"))

	router := transporthttp.NewRouter(handler, authHandler, authSvc)
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.WithField("addr", addr).Info("starting API server")

	if err := router.Run(addr); err != nil {
		log.WithError(err).Fatal("server stopped")
	}
}
