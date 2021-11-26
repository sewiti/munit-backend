package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/auth"
	"github.com/sewiti/munit-backend/internal/model"
	"github.com/sewiti/munit-backend/pkg/id"
)

type contextKey int

const (
	// gorilla/mux uses 0 and 1
	userKey contextKey = 2
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 {
			respondUnauthorized(w)
			return
		}

		var uid id.ID
		switch authParts[0] {
		case "Bearer":
			subject, err := auth.VerifyJWT(authParts[1])
			if err != nil {
				log.WithError(err).Debug("unable to verify jwt")
				respondUnauthorized(w)
				return
			}
			uid = id.ID(subject)
		default:
			respondUnauthorized(w)
			return
		}

		_, err := model.GetUser(r.Context(), uid)
		if err != nil {
			respondUnauthorized(w)
			return
		}

		if ids, err := getIDs(r, projectID); err == nil && len(ids) == 1 {
			err = verifyProjectAssociate(r.Context(), ids[0], uid)
			if err != nil {
				respondErr(w, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), userKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyProjectAssociate(ctx context.Context, project, user id.ID) error {
	p, err := model.GetProject(ctx, project)
	if err != nil {
		return err
	}
	if user == p.Owner {
		return nil
	}

	for _, contrib := range p.Contributors {
		if user == contrib {
			return nil
		}
	}
	return errForbidden
	// return model.ErrNotFound // Fake 404
}

func getUser(r *http.Request) (id.ID, error) {
	var uid id.ID
	if v := r.Context().Value(userKey); v != nil {
		uid = v.(id.ID)
	}
	return uid, uid.Validate()
}
