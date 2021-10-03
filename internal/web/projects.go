package web

import (
	"net/http"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/sewiti/munit-backend/internal/model"
)

func projectGetAll(w http.ResponseWriter, r *http.Request) {
	p, err := model.GetAllProjects()
	if err != nil {
		respondNotFound(w, err)
		return
	}
	respondOK(w, p)
}

func projectGet(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	p, err := model.GetProject(uuids[0])
	if err != nil {
		respondNotFound(w, err)
		return
	}
	respondOK(w, p)
}

func projectPost(w http.ResponseWriter, r *http.Request) {
	var p model.Project
	if err := decodeJSON(r, &p); err != nil {
		respondBadRequest(w, err)
		return
	}

	var err error
	p.UUID, err = uuid.NewRandom()
	if err != nil {
		log.WithError(err).Error("unable to generate uuid")
		respondInternalError(w)
		return
	}

	if err = model.AddProject(&p); err != nil {
		respondBadRequest(w, err)
		return
	}
	respond(w, p, http.StatusCreated)
}

func projectPut(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	var p model.Project
	if err := decodeJSON(r, &p); err != nil {
		respondBadRequest(w, err)
		return
	}
	p.UUID = uuids[0] // Disallow changing UUID

	if err = model.EditProject(&p); err != nil {
		respondBadRequest(w, err)
		return
	}
	respondOK(w, p)
}

func projectDelete(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	if err := model.DeleteProject(uuids[0]); err != nil {
		respondNotFound(w, err)
		return
	}
	respondNoContent(w)
}
