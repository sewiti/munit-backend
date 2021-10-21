package web

import (
	"net/http"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/id"
	"github.com/sewiti/munit-backend/internal/model"
)

func projectGetAll(w http.ResponseWriter, r *http.Request) {
	p, err := model.GetAllProjects()
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

	p, err := model.GetProject(prID)
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

	if err = model.AddProject(&p); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, p, http.StatusCreated)
}

func projectPut(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	var p model.Project
	if err := decodeJSON(r, &p); err != nil {
		respondErr(w, err)
		return
	}
	p.ID = prID // Disallow changing ID

	if err = model.EditProject(&p); err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, p)
}

func projectDelete(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	if err := model.DeleteProject(prID); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
