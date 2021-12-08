package web

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sewiti/munit-backend/internal/config"
)

const (
	userID    = "userID"    // User ID path key
	projectID = "projectID" // Project ID path key
	commitID  = "commitID"  // Commit ID path key
	fileID    = "fileID"    // File ID path key

	idPattern = "[A-Za-z0-9]+"

	defaultBodyLimit     = 1024 * 1024      // 1MiB
	defaultFileBodyLimit = 1024 * 1024 * 50 // 50MiB
)

func NewRouter(cfg *config.Munit) http.Handler {
	const (
		userVar    = "{" + userID + ":" + idPattern + "}"
		projectVar = "{" + projectID + ":" + idPattern + "}"
		commitVar  = "{" + commitID + ":" + idPattern + "}"
		fileVar    = "{" + fileID + ":" + idPattern + "}"
	)
	r := mux.NewRouter()

	// Auth
	r.Methods("POST").Path("/register").HandlerFunc(registerPost)
	r.Methods("POST").Path("/login").HandlerFunc(loginPost)

	// Profile
	profile := r.PathPrefix("/profile").Subrouter()
	profile.Use(authMiddleware)
	profile.Methods("GET").Path("/" + userVar).HandlerFunc(profileGet)
	profile.Methods("GET").Path("").HandlerFunc(profileSelfGet)
	profile.Methods("PATCH").Path("").HandlerFunc(profilePatch)
	profile.Methods("DELETE").Path("").HandlerFunc(profileDelete)

	// Project
	project := r.PathPrefix("/projects").Subrouter()
	project.Use(authMiddleware)
	project.Methods("GET").Path("").HandlerFunc(projectGetAll)
	project.Methods("POST").Path("").HandlerFunc(projectPost)
	project.Methods("GET").Path("/" + projectVar).HandlerFunc(projectGet)
	project.Methods("PATCH").Path("/" + projectVar).HandlerFunc(projectPatch)
	project.Methods("DELETE").Path("/" + projectVar).HandlerFunc(projectDelete)

	// Commit
	commit := project.PathPrefix("/" + projectVar + "/commits").Subrouter()
	commit.Methods("GET").Path("").HandlerFunc(commitGetAll)
	commit.Methods("POST").Path("").HandlerFunc(commitPost)
	commit.Methods("GET").Path("/" + commitVar).HandlerFunc(commitGet)
	commit.Methods("PATCH").Path("/" + commitVar).HandlerFunc(commitPatch)
	commit.Methods("DELETE").Path("/" + commitVar).HandlerFunc(commitDelete)

	// File
	file := commit.PathPrefix("/" + commitVar + "/files").Subrouter()
	file.Methods("GET").Path("").HandlerFunc(fileGetAll)
	file.Methods("POST").Path("").HandlerFunc(filePost)
	file.Methods("GET").Path("/" + fileVar).HandlerFunc(fileGet)
	file.Methods("PATCH").Path("/" + fileVar).HandlerFunc(filePatch)
	file.Methods("DELETE").Path("/" + fileVar).HandlerFunc(fileDelete)

	// Setup CORS
	origins := handlers.AllowedOrigins([]string{cfg.AllowedOrigin})
	headers := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"})
	return handlers.CORS(origins, headers, methods)(r)
}
