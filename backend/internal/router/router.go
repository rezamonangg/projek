package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/admin"
	"github.com/monachy/projek/internal/auth"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/file"
	"github.com/monachy/projek/internal/member"
	"github.com/monachy/projek/internal/project"
	"github.com/monachy/projek/internal/task"
	"github.com/monachy/projek/internal/wiki"
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
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.Handler(corsOpts))

	memberRepo := member.NewRepository(db)
	memberService := member.NewService(memberRepo)
	memberHandler := member.NewHandler(memberService)

	sessionStore := auth.NewSessionStore(redis, 7*24*time.Hour)
	authService := auth.NewAuthService(memberRepo, sessionStore)
	r.Use(auth.AuthMiddleware(authService))

	authHandler := auth.NewHandler(authService, memberService)

	r.Get("/health", healthHandler)
	r.Get("/ready", readinessHandler(db, redis))
	r.Get("/health/detailed", detailedHealthHandler(db, redis))

	r.Mount("/auth", authHandler.Routes())
	r.Mount("/members", memberHandler.Routes())

	projectRepo := project.NewRepository(db)
	projectEpicRepo := project.NewEpicRepository(db)
	projectBoardRepo := project.NewBoardRepository(db)
	projectService := project.NewService(projectRepo, projectEpicRepo, projectBoardRepo, nil, nil)
	projectHandler := project.NewHandler(projectService)
	r.Mount("/projects", projectHandler.Routes())

	epicHandler := project.NewEpicHandler(projectService)
	boardHandler := project.NewBoardHandler(projectService)

	taskRepo := task.NewRepository(db)
	taskLabelRepo := task.NewLabelRepository(db)
	taskService := task.NewService(taskRepo, taskLabelRepo)
	taskHandler := task.NewHandler(taskService)

	taskFullService := task.NewTaskService(taskRepo, taskLabelRepo)
	labelHandler := task.NewLabelHandler(taskFullService)

	wikiRepo := wiki.NewRepository(db)
	wikiService := wiki.NewService(wikiRepo)
	wikiHandler := wiki.NewHandler(wikiService)

	fileRepo := file.NewRepository(db)
	fileStorage := file.NewLocalStorage("/var/lib/projek/uploads", "/uploads")
	fileService := file.NewService(fileRepo, fileStorage)
	fileHandler := file.NewHandler(fileService)

	adminSettingsRepo := admin.NewSettingsRepository(db)
	adminStatsRepo := admin.NewStatsRepository(db)
	adminService := admin.NewService(adminSettingsRepo, adminStatsRepo)
	adminHandler := admin.NewHandler(adminService)

	r.Route("/projects/{projectId}", func(r chi.Router) {
		r.Mount("/epics", epicHandler.ProjectRoutes())
		r.Mount("/boards", boardHandler.ProjectRoutes())
		r.Mount("/labels", labelHandler.ProjectRoutes())
		r.Mount("/wiki", wikiHandler.ProjectRoutes())
		r.Mount("/files", fileHandler.ProjectRoutes())
	})
	r.Mount("/epics", epicHandler.Routes())
	r.Mount("/boards", boardHandler.Routes())
	r.Mount("/wiki", wikiHandler.Routes())
	r.Mount("/files", fileHandler.Routes())

	r.Route("/boards/{boardId}", func(r chi.Router) {
		r.Mount("/tasks", taskHandler.BoardRoutes())
	})
	r.Mount("/tasks", taskHandler.Routes())

	r.Route("/tasks/{taskId}", func(r chi.Router) {
		r.Mount("/labels", labelHandler.TaskRoutes())
	})

	r.Mount("/admin", adminHandler.Routes())

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
