package web

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sewiti/munit-backend/internal/config"
)

const (
	projectUUID = "projectUUID" // Project UUID path key
	commitUUID  = "commitUUID"  // Commit UUID path key
	fileUUID    = "fileUUID"    // File UUID path key

	bodyLimit     = 1024 * 1024      // 1MiB
	fileBodyLimit = 1024 * 1024 * 50 // 50MiB
)

func NewRouter(cfg *config.Munit) http.Handler {
	r := mux.NewRouter()

	// Setup CORS
	origins := handlers.AllowedOrigins([]string{cfg.AllowedOrigin})
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost,
		http.MethodPut, http.MethodDelete})

	return handlers.CORS(origins, headers, methods)(r)
}
