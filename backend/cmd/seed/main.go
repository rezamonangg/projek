package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/scripts"
)

func main() {
	if os.Getenv("ENV") != "development" {
		log.Fatal("seed can only run in development environment. Set ENV=development")
	}

	cfg := common.LoadConfig()

	if err := common.Validate(cfg); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	logger := common.NewLogger(cfg)
	logger.Info().Msg("starting seed process")

	db, err := common.NewDatabase(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := scripts.Seed(ctx, db, logger); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	fmt.Println("\n✓ Seed completed successfully!")
	fmt.Println("\nDummy data created:")
	fmt.Println("  - Community: Test Community (slug: test-community)")
	fmt.Println("  - Admin: admin@example.com / password: admin123")
}
