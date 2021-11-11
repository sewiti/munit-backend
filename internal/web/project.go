package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/model"
	"github.com/sewiti/munit-backend/pkg/id"
)

func projectGetAll(w http.ResponseWriter, r *http.Request) {
	uid, err := getUser(r)
	if err != nil {
		log.WithError(err).Error("unable to get user from context")
		respondInternalError(w)
		return
	}

	p, err := model.GetAllProjects(r.Context(), uid)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, p)
}

func projectGet(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID)
	if err != nil {
		respondErr(w, err)
		return
	}

	p, err := model.GetProject(r.Context(), ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, p)
}

func projectPost(w http.ResponseWriter, r *http.Request) {
	uid, err := getUser(r)
	if err != nil {
		log.WithError(err).Error("unable to get user from context")
		respondInternalError(w)
		return
	}

	var p model.Project
	if err := decodeJSON(r, &p); err != nil {
		respondErr(w, err)
		return
	}

	p.ID, err = id.New()
	if err != nil {
		log.WithError(err).Error("unable to generate id")
		respondInternalError(w)
		return
	}
	now := time.Now().Truncate(time.Second)
	p.Owner = uid
	p.Created = now
	p.Modified = now

	if err = model.InsertProject(r.Context(), &p); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, p, http.StatusCreated)
}

func projectPatch(w http.ResponseWriter, r *http.Request) {
	if err := assertJSON(r); err != nil {
		respondErr(w, err)
		return
	}
	data, err := ioutil.ReadAll(io.LimitReader(r.Body, defaultBodyLimit))
	if err != nil {
		respondErr(w, err)
		return
	}

	ids, err := getIDs(r, projectID)
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

	p, err := model.UpdateProject(r.Context(), ids[0], func(p *model.Project) error {
		if p.Owner != uid {
			return errForbidden
		}

		orig := *p
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		p.ID = orig.ID
		p.Created = orig.Created
		p.Modified = time.Now().Truncate(time.Second)
		if p.Contributors == nil {
			p.Contributors = make([]id.ID, 0)
		}
		return nil
	})
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, p)
}

func projectDelete(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID)
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

	p, err := model.GetProject(r.Context(), ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	if p.Owner != uid {
		respondErr(w, errForbidden)
		return
	}

	if err = model.DeleteProject(r.Context(), ids[0]); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
