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

	// Projects
	pr := r.PathPrefix("/projects").Subrouter()
	pr.Methods("GET").Path("").HandlerFunc(projectGetAll)
	pr.Methods("POST").Path("").HandlerFunc(projectPost)
	pr.Methods("GET").Path("/{" + projectUUID + "}").HandlerFunc(projectGet)
	pr.Methods("PUT").Path("/{" + projectUUID + "}").HandlerFunc(projectPut)
	pr.Methods("DELETE").Path("/{" + projectUUID + "}").HandlerFunc(projectDelete)

	// Commits
	cr := pr.PathPrefix("/{" + projectUUID + "}/commits").Subrouter()
	cr.Methods("GET").Path("").HandlerFunc(commitGetAll)
	cr.Methods("POST").Path("").HandlerFunc(commitPost)
	cr.Methods("GET").Path("/{" + commitUUID + "}").HandlerFunc(commitGet)
	cr.Methods("PUT").Path("/{" + commitUUID + "}").HandlerFunc(commitPut)
	cr.Methods("DELETE").Path("/{" + commitUUID + "}").HandlerFunc(commitDelete)

	// Setup CORS
	origins := handlers.AllowedOrigins([]string{cfg.AllowedOrigin})
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost,
		http.MethodPut, http.MethodDelete})

	return handlers.CORS(origins, headers, methods)(r)
}
