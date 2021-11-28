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

func commitGetAll(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID)
	if err != nil {
		respondErr(w, err)
		return
	}
	_, err = model.GetProject(r.Context(), ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	c, err := model.GetAllCommits(r.Context(), ids[0])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, c)
}

func commitGet(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	c, err := model.GetCommit(r.Context(), ids[0], ids[1])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, c)
}

func commitPost(w http.ResponseWriter, r *http.Request) {
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

	var c model.Commit
	if err = decodeJSON(r, &c); err != nil {
		respondErr(w, err)
		return
	}

	c.ID, err = id.New()
	if err != nil {
		log.WithError(err).Error("unable to make id")
		respondInternalError(w)
		return
	}
	now := time.Now().Truncate(time.Second)
	c.Created = now
	c.Modified = now
	c.User = uid
	c.Project = ids[0]

	if err = model.InsertCommit(r.Context(), &c); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, c, http.StatusCreated)
}

func commitPatch(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}
	if err := assertJSON(r); err != nil {
		respondErr(w, err)
		return
	}
	data, err := ioutil.ReadAll(io.LimitReader(r.Body, defaultBodyLimit))
	if err != nil {
		respondErr(w, err)
		return
	}

	c, err := model.EditCommit(r.Context(), ids[0], ids[1], func(c *model.Commit) error {
		orig := *c
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
		c.ID = orig.ID
		c.Created = orig.Created
		c.Modified = time.Now().Truncate(time.Second)
		c.Project = orig.Project
		c.User = orig.User
		return nil
	})
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, c)
}

func commitDelete(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}
	err = model.DeleteCommit(r.Context(), ids[0], ids[1])
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
