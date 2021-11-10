package web

import (
	"crypto/rand"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/auth"
	"github.com/sewiti/munit-backend/internal/model"
	"github.com/sewiti/munit-backend/pkg/id"
)

func register(w http.ResponseWriter, r *http.Request) {
	u := new(model.User)
	if err := decodeJSON(r, u); err != nil {
		respondErr(w, err)
		return
	}

	var err error
	u.ID, err = id.New()
	if err != nil {
		respondErr(w, err)
		return
	}
	u.Salt, err = auth.MakeSalt(rand.Reader)
	if err != nil {
		respondErr(w, err)
		return
	}
	u.Hash = auth.HashPasswd([]byte(u.Password), u.Salt)
	u.Password = "" // not needed anymore
	u.Created = time.Now()
	u.Modified = time.Now()

	if err = model.InsertUser(r.Context(), u); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, u, http.StatusCreated)
}

func login(w http.ResponseWriter, r *http.Request) {
	u := new(model.User)
	if err := decodeJSON(r, u); err != nil {
		respondErr(w, err)
		return
	}

	if u.Email == "" {
		respondMsg(w, "email is empty", 400)
		return
	}
	if u.Password == "" {
		respondMsg(w, "password is empty", 400)
		return
	}

	dbUsr, err := model.GetUserByEmail(r.Context(), u.Email)
	if err != nil {
		respondMsg(w, "Unauthorized", 401)
		return
	}

	if !auth.VerifyPasswd(dbUsr.Hash, []byte(u.Password), dbUsr.Salt) {
		respondMsg(w, "Unauthorized", 401)
		return
	}
	token, err := auth.MakeJWT(string(dbUsr.ID))
	if err != nil {
		log.WithError(err).WithField("user", dbUsr.ID).Error("unable to make jwt")
		respondInternalError(w)
		return
	}
	respondOK(w, struct {
		Token string `json:"token"`
	}{token})
}

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
			respondErr(w, err)
			return
		}

		// prID, _, err := getIDs(r)
		// if err != nil {
		// 	respondErr(w, err)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
