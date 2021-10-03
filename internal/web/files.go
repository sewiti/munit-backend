package web

import (
	"net/http"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/sewiti/munit-backend/internal/model"
)

func fileGetAll(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	f, err := model.GetAllFiles(uuids[1])
	if err != nil {
		respondNotFound(w, err)
		return
	}
	respondOK(w, f)
}

func fileGet(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID, fileUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	f, err := model.GetFile(uuids[1], uuids[2])
	if err != nil {
		respondNotFound(w, err)
		return
	}
	respondOK(w, f)
}

func filePost(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	var f model.File
	if err = decodeJSONLimit(r, &f, fileBodyLimit); err != nil {
		respondBadRequest(w, err)
		return
	}

	f.UUID, err = uuid.NewRandom()
	if err != nil {
		log.WithError(err).Error("unable to generate uuid")
		respondInternalError(w)
		return
	}
	f.Commit = uuids[1]

	if err = model.AddFile(&f); err != nil {
		respondBadRequest(w, err)
		return
	}
	respond(w, f, http.StatusCreated)
}

func filePut(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID, fileUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	var f model.File
	if err = decodeJSONLimit(r, &f, fileBodyLimit); err != nil {
		respondBadRequest(w, err)
		return
	}

	f.Commit = uuids[1] // Disallow changing Project
	f.UUID = uuids[2]   // Disallow changing UUID

	if err = model.EditFile(&f); err != nil {
		respondBadRequest(w, err)
		return
	}
	respondOK(w, f)
}

func fileDelete(w http.ResponseWriter, r *http.Request) {
	uuids, err := getUUIDs(r, projectUUID, commitUUID, fileUUID)
	if err != nil {
		respondBadRequest(w, err)
		return
	}

	if err := model.DeleteFile(uuids[1], uuids[2]); err != nil {
		respondNotFound(w, err)
		return
	}
	respondNoContent(w)
}
