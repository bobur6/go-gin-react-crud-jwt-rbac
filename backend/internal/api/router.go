package api

import (
	"strings"
	"time"

	"assignment3/backend/internal/auth"
	"assignment3/backend/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all routes and middleware.
func SetupRouter(store *store.Store, jwtService *auth.JWTService, allowedOrigins []string, allowAll bool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	_ = router.SetTrustedProxies(nil)

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if allowAll {
		corsConfig.AllowAllOrigins = true
	} else {
		originSet := make(map[string]struct{}, len(allowedOrigins)*2)
		for _, origin := range allowedOrigins {
			norm := normalizeOrigin(origin)
			if norm == "" {
				continue
			}
			originSet[norm] = struct{}{}
		}

		corsConfig.AllowOriginFunc = func(origin string) bool {
			norm := normalizeOrigin(origin)
			if norm == "" {
				return false
			}
			_, ok := originSet[norm]
			return ok
		}
	}

	router.Use(cors.New(corsConfig))

	handler := NewHandler(store, jwtService)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/health", handler.Health)
		apiGroup.POST("/register", handler.Register)
		apiGroup.POST("/login", handler.Login)

		items := apiGroup.Group("/items")
		items.Use(auth.AuthMiddleware(jwtService))
		{
			items.GET("", handler.ListItems)
			items.POST("", handler.CreateItem)
			items.GET("/:id", handler.GetItem)
			items.PUT("/:id", handler.UpdateItem)
			items.DELETE("/:id", auth.RequireRoles("admin"), handler.DeleteItem)
		}

		users := apiGroup.Group("/users")
		users.Use(auth.AuthMiddleware(jwtService), auth.RequireRoles("admin"))
		{
			users.GET("", handler.ListUsers)
			users.DELETE("/:id", handler.DeleteUser)
		}
	}

	return router
}

func normalizeOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return ""
	}
	origin = strings.TrimRight(origin, "/")
	return strings.ToLower(origin)
}
