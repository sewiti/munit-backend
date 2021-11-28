package web

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/auth"
	"github.com/sewiti/munit-backend/internal/model"
	"github.com/sewiti/munit-backend/pkg/id"
)

func registerPost(w http.ResponseWriter, r *http.Request) {
	u := new(model.User)
	if err := decodeJSON(r, u); err != nil {
		respondErr(w, err)
		return
	}

	var err error
	u.ID, err = id.New()
	if err != nil {
		log.WithError(err).Error("unable to make new id")
		respondInternalError(w)
		return
	}
	u.Salt, err = auth.MakeSalt(rand.Reader)
	if err != nil {
		log.WithError(err).Error("unable to make salt")
		respondInternalError(w)
		return
	}
	u.PasswdHash = auth.HashPasswd([]byte(u.Password), u.Salt)
	now := time.Now().Truncate(time.Second)
	u.Created = now
	u.Modified = now

	if err = model.InsertUser(r.Context(), u); err != nil {
		respondErr(w, err)
		return
	}
	u.Password = "" // never ouput it
	respond(w, u, http.StatusCreated)
}

func loginGet(w http.ResponseWriter, r *http.Request) {
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
		respondUnauthorized(w)
		return
	}

	if !auth.VerifyPasswd(dbUsr.PasswdHash, []byte(u.Password), dbUsr.Salt) {
		respondUnauthorized(w)
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

func profileGet(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, userID)
	if err != nil {
		respondErr(w, err)
		return
	}
	u, err := model.GetUser(r.Context(), ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, u)
}

func profileSelfGet(w http.ResponseWriter, r *http.Request) {
	uid, err := getUser(r)
	if err != nil {
		log.WithError(err).Error("unable to get user from context")
		respondInternalError(w)
		return
	}
	u, err := model.GetUser(r.Context(), uid)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, u)
}

func profilePatch(w http.ResponseWriter, r *http.Request) {
	if err := assertJSON(r); err != nil {
		respondErr(w, err)
		return
	}
	data, err := ioutil.ReadAll(io.LimitReader(r.Body, defaultBodyLimit))
	if err != nil {
		respondErr(w, err)
		return
	}

	uid, err := getUser(r)
	if err != nil {
		log.WithError(err).Error("unable to get user from context")
		respondInternalError(w)
		return
	}

	u, err := model.UpdateUser(r.Context(), uid, func(u *model.User) error {
		orig := u.Copy()
		if err := json.Unmarshal(data, u); err != nil {
			return err
		}
		u.ID = orig.ID
		u.Created = orig.Created
		u.Modified = time.Now().Truncate(time.Second)

		if u.Password == "" { // means we are not changing password
			u.PasswdHash = orig.PasswdHash
			u.Salt = orig.Salt
		} else {
			u.Salt, err = auth.MakeSalt(rand.Reader)
			if err != nil {
				log.WithError(err).Error("unable to make salt")
				return errInternalError
			}
			u.PasswdHash = auth.HashPasswd([]byte(u.Password), u.Salt)
		}
		return nil
	})
	if err != nil {
		respondErr(w, err)
		return
	}
	u.Password = "" // never ouput it
	respondOK(w, u)
}

func profileDelete(w http.ResponseWriter, r *http.Request) {
	uid, err := getUser(r)
	if err != nil {
		log.WithError(err).Error("unable to get user from context")
		respondInternalError(w)
		return
	}
	if err = model.DeleteUser(r.Context(), uid); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
