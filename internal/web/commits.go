package web

import (
	"net/http"

	"github.com/sewiti/munit-backend/internal/model"
)

func commitGetAll(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	c, err := model.GetAllCommits(prID)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, c)
}

func commitGet(w http.ResponseWriter, r *http.Request) {
	prID, ids, err := getIDs(r, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	c, err := model.GetCommit(prID, ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, c)
}

func commitPost(w http.ResponseWriter, r *http.Request) {
	prID, _, err := getIDs(r)
	if err != nil {
		respondErr(w, err)
		return
	}

	var c model.Commit
	if err = decodeJSON(r, &c); err != nil {
		respondErr(w, err)
		return
	}

	c.Project = prID

	if err = model.AddCommit(&c); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, c, http.StatusCreated)
}

func commitPut(w http.ResponseWriter, r *http.Request) {
	prID, ids, err := getIDs(r, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	var c model.Commit
	if err = decodeJSON(r, &c); err != nil {
		respondErr(w, err)
		return
	}

	c.Project = prID // Disallow changing Project
	c.ID = ids[0]    // Disallow changing ID

	if err = model.EditCommit(&c); err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, c)
}

func commitDelete(w http.ResponseWriter, r *http.Request) {
	prID, ids, err := getIDs(r, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	if err := model.DeleteCommit(prID, ids[0]); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
