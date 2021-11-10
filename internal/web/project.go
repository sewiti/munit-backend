package web

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/model"
	"github.com/sewiti/munit-backend/pkg/id"
)

func projectGetAll(w http.ResponseWriter, r *http.Request) {
	p, err := model.GetAllProjects(r.Context())
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, p)
}

func projectGet(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	p, err := model.GetProject(r.Context(), prID)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, p)
}

func projectPost(w http.ResponseWriter, r *http.Request) {
	var p model.Project
	if err := decodeJSON(r, &p); err != nil {
		respondErr(w, err)
		return
	}

	var err error
	p.ID, err = id.New()
	if err != nil {
		log.WithError(err).Error("unable to generate uuid")
		respondInternalError(w)
		return
	}
	now := time.Now()
	p.Created = now
	p.Modified = now

	if err = model.InsertProject(r.Context(), &p); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, p, http.StatusCreated)
}

func projectPatch(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	var ret *model.Project
	err = model.UpdateProject(r.Context(), prID, func(p *model.Project) error {
		orig := *p
		if err := decodeJSON(r, &p); err != nil {
			return err
		}
		p.ID = orig.ID
		p.Created = orig.Created
		p.Modified = time.Now()
		if p.Contributors == nil {
			p.Contributors = make([]id.ID, 0)
		}
		ret = p
		return nil
	})
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, ret)
}

func projectDelete(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	if err := model.DeleteProject(r.Context(), prID); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
