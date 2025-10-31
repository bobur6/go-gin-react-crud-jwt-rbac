package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"assignment3/backend/internal/api"
	"assignment3/backend/internal/auth"
	"assignment3/backend/internal/store"
)

func main() {
	port := getenvDefault("PORT", "8080")
	jwtSecret := getenvDefault("JWT_SECRET", "change-me-in-production")
	jwtIssuer := getenvDefault("JWT_ISSUER", "assignment3-backend")
	expiryMinutes := getenvIntDefault("JWT_EXPIRY_MINUTES", 60)

	st := store.NewStore()

	adminUsername := getenvDefault("ADMIN_USERNAME", "admin")
	adminPassword := getenvDefault("ADMIN_PASSWORD", "admin123")

	if _, created, err := st.EnsureAdminUser(adminUsername, adminPassword); err != nil {
		log.Fatalf("failed to ensure admin user: %v", err)
	} else if created {
		log.Printf("created default admin user '%s'", adminUsername)
	} else {
		log.Printf("admin user '%s' already exists", adminUsername)
	}

	// Seed with an example item to illustrate API responses.
	if _, err := st.CreateItem(adminUsername, "Welcome Item", "You can edit or delete this item from the React app."); err != nil {
		log.Printf("warning: failed to seed welcome item: %v", err)
	}

	jwtService := auth.NewJWTService(jwtSecret, jwtIssuer, time.Duration(expiryMinutes)*time.Minute)
	origins, allowAll := loadAllowedOrigins(port)

	router := api.SetupRouter(st, jwtService, origins, allowAll)

	log.Printf("server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvIntDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return fallback
}

func parseCSVEnv(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		normalized := strings.TrimSpace(part)
		if normalized == "" {
			continue
		}
		normalized = strings.TrimRight(normalized, "/")
		lower := strings.ToLower(normalized)
		if _, exists := seen[lower]; exists {
			continue
		}
		seen[lower] = struct{}{}
		result = append(result, lower)
	}
	return result
}

func loadAllowedOrigins(port string) ([]string, bool) {
	raw := getenvDefault("FRONTEND_ORIGINS", "")
	origins := parseCSVEnv(raw)

	for _, origin := range origins {
		if origin == "*" {
			return nil, true
		}
	}

	defaultOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		fmt.Sprintf("http://localhost:%s", port),
		fmt.Sprintf("http://127.0.0.1:%s", port),
	}

	originSet := make(map[string]struct{}, len(origins)+len(defaultOrigins))
	for _, origin := range origins {
		originSet[origin] = struct{}{}
	}

	for _, origin := range defaultOrigins {
		norm := strings.ToLower(strings.TrimRight(strings.TrimSpace(origin), "/"))
		if _, exists := originSet[norm]; !exists {
			origins = append(origins, norm)
			originSet[norm] = struct{}{}
		}
	}

	return origins, len(origins) == 0
}
