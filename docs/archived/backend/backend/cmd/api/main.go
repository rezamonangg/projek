package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/router"
)

func main() {
	cfg := common.LoadConfig()

	if err := common.Validate(cfg); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	logger := common.NewLogger(cfg)
	logger.Info().Str("version", "1.0.0").Msg("starting projek backend")

	db, err := common.NewDatabase(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redis, err := common.NewRedis(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	r := router.NewRouter(cfg, logger, db, redis)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info().Str("addr", addr).Msg("starting HTTP server")

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal().Err(err).Msg("HTTP server failed")
	}
}
