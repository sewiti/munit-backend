package web

import (
	"net/http"

	"github.com/sewiti/munit-backend/internal/model"
)

func fileGetAll(w http.ResponseWriter, r *http.Request) {
	_, ids, err := getIDs(r, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	f, err := model.GetAllFiles(ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, f)
}

func fileGet(w http.ResponseWriter, r *http.Request) {
	_, ids, err := getIDs(r, commitID, fileID)
	if err != nil {
		respondErr(w, err)
		return
	}

	f, err := model.GetFile(ids[0], ids[1])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, f)
}

func filePost(w http.ResponseWriter, r *http.Request) {
	_, ids, err := getIDs(r, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	var f model.File
	if err = decodeJSONLimit(r, &f, fileBodyLimit); err != nil {
		respondErr(w, err)
		return
	}

	f.Commit = ids[0]

	if err = model.AddFile(&f); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, f, http.StatusCreated)
}

func filePut(w http.ResponseWriter, r *http.Request) {
	_, ids, err := getIDs(r, commitID, fileID)
	if err != nil {
		respondErr(w, err)
		return
	}

	var f model.File
	if err = decodeJSONLimit(r, &f, fileBodyLimit); err != nil {
		respondErr(w, err)
		return
	}

	f.Commit = ids[0] // Disallow changing Project
	f.ID = ids[1]     // Disallow changing ID

	if err = model.EditFile(&f); err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, f)
}

func fileDelete(w http.ResponseWriter, r *http.Request) {
	_, ids, err := getIDs(r, commitID, fileID)
	if err != nil {
		respondErr(w, err)
		return
	}

	if err := model.DeleteFile(ids[0], ids[1]); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
