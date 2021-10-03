package web

import (
	"net/http"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/sewiti/munit-backend/internal/model"
)

func commitGetAll(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	c, err := model.GetAllCommits(uuids[0])
	if err != nil {
		respondNotFound(w, err)
		return
	}
	respondOK(w, c)
}

func commitGet(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	c, err := model.GetCommit(uuids[0], uuids[1])
	if err != nil {
		respondNotFound(w, err)
		return
	}
	respondOK(w, c)
}

func commitPost(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	var c model.Commit
	if err = decodeJSON(r, &c); err != nil {
		respondBadRequest(w, err)
		return
	}

	c.UUID, err = uuid.NewRandom()
	if err != nil {
		log.WithError(err).Error("unable to generate uuid")
		respondInternalError(w)
		return
	}
	c.Project = uuids[0]

	if err = model.AddCommit(&c); err != nil {
		respondBadRequest(w, err)
		return
	}
	respond(w, c, http.StatusCreated)
}

func commitPut(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	var c model.Commit
	if err = decodeJSON(r, &c); err != nil {
		respondBadRequest(w, err)
		return
	}

	c.Project = uuids[0] // Disallow changing Project
	c.UUID = uuids[1]    // Disallow changing UUID

	if err = model.EditCommit(&c); err != nil {
		respondBadRequest(w, err)
		return
	}
	respondOK(w, c)
}

func commitDelete(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	if err := model.DeleteCommit(uuids[0], uuids[1]); err != nil {
		respondNotFound(w, err)
		return
	}
	respondNoContent(w)
}
