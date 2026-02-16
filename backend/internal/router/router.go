package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/common"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRouter(cfg *common.Config, logger zerolog.Logger, db *pgxpool.Pool, redis *redis.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(common.RequestLogger(logger))

	corsOpts := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.Handler(corsOpts))

	r.Get("/health", healthHandler)
	r.Get("/ready", readinessHandler(db, redis))
	r.Get("/health/detailed", detailedHealthHandler(db, redis))

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readinessHandler(db *pgxpool.Pool, redis *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := db.Ping(ctx); err != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}
		if err := redis.Ping(ctx).Err(); err != nil {
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func detailedHealthHandler(db *pgxpool.Pool, redis *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		status := map[string]interface{}{
			"status": "healthy",
			"checks": map[string]interface{}{},
		}

		if err := db.Ping(ctx); err != nil {
			status["status"] = "unhealthy"
			status["checks"].(map[string]interface{})["database"] = map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			status["checks"].(map[string]interface{})["database"] = map[string]string{
				"status": "healthy",
			}
		}

		if err := redis.Ping(ctx).Err(); err != nil {
			status["status"] = "unhealthy"
			status["checks"].(map[string]interface{})["redis"] = map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			status["checks"].(map[string]interface{})["redis"] = map[string]string{
				"status": "healthy",
			}
		}

		if status["status"] == "unhealthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		common.WriteJSON(w, http.StatusOK, status)
	}
}
