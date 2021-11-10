package web

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sewiti/munit-backend/internal/config"
)

const (
	projectID = "projectID" // Project ID path key
	commitID  = "commitID"  // Commit ID path key
	fileID    = "fileID"    // File ID path key

	projectPattern = "[A-Za-z0-9]+"
	commitPattern  = "[0-9]+"
	filePattern    = "[0-9]+"

	bodyLimit     = 1024 * 1024      // 1MiB
	fileBodyLimit = 1024 * 1024 * 50 // 50MiB
)

func NewRouter(cfg *config.Munit) http.Handler {
	const (
		projectVar = "{" + projectID + ":" + projectPattern + "}"
		commitVar  = "{" + commitID + ":" + commitPattern + "}"
		fileVar    = "{" + fileID + ":" + filePattern + "}"
	)
	r := mux.NewRouter()

	// Auth
	r.Methods("POST").Path("/register").HandlerFunc(register)
	r.Methods("GET").Path("/login").HandlerFunc(login)

	// Projects
	pr := r.PathPrefix("/projects").Subrouter()
	pr.Use(authMiddleware)
	pr.Methods("GET").Path("").HandlerFunc(projectGetAll)
	pr.Methods("POST").Path("").HandlerFunc(projectPost)
	pr.Methods("GET").Path("/" + projectVar).HandlerFunc(projectGet)
	pr.Methods("PATCH").Path("/" + projectVar).HandlerFunc(projectPatch)
	pr.Methods("DELETE").Path("/" + projectVar).HandlerFunc(projectDelete)

	// Commits
	cr := pr.PathPrefix("/" + projectVar + "/commits").Subrouter()
	cr.Methods("GET").Path("").HandlerFunc(commitGetAll)
	cr.Methods("POST").Path("").HandlerFunc(commitPost)
	cr.Methods("GET").Path("/" + commitVar).HandlerFunc(commitGet)
	cr.Methods("PUT").Path("/" + commitVar).HandlerFunc(commitPut)
	cr.Methods("DELETE").Path("/" + commitVar).HandlerFunc(commitDelete)

	// Files
	fr := cr.PathPrefix("/" + commitVar + "/files").Subrouter()
	fr.Methods("GET").Path("").HandlerFunc(fileGetAll)
	fr.Methods("POST").Path("").HandlerFunc(filePost)
	fr.Methods("GET").Path("/" + fileVar).HandlerFunc(fileGet)
	fr.Methods("PUT").Path("/" + fileVar).HandlerFunc(filePut)
	fr.Methods("DELETE").Path("/" + fileVar).HandlerFunc(fileDelete)

	// Setup CORS
	origins := handlers.AllowedOrigins([]string{cfg.AllowedOrigin})
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"})
	return handlers.CORS(origins, headers, methods)(r)
}
