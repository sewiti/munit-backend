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

func fileGetAll(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}
	_, err = model.GetCommit(r.Context(), ids[0], ids[1])
	if err != nil {
		respondErr(w, err)
		return
	}

	f, err := model.GetAllFiles(r.Context(), ids[0], ids[1])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, f)
}

func fileGet(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID, fileID)
	if err != nil {
		respondErr(w, err)
		return
	}

	f, err := model.GetFile(r.Context(), ids[0], ids[1], ids[2])
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, f)
}

func filePost(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID)
	if err != nil {
		respondErr(w, err)
		return
	}

	var f model.File
	if err = decodeJSONLimit(r, &f, defaultFileBodyLimit); err != nil {
		respondErr(w, err)
		return
	}

	f.ID, err = id.New()
	if err != nil {
		log.WithError(err).Error("unable to make id")
		respondInternalError(w)
		return
	}
	now := time.Now().Truncate(time.Second)
	f.Created = now
	f.Modified = now
	f.Project = ids[0]
	f.Commit = ids[1]

	if err = model.InsertFile(r.Context(), &f); err != nil {
		respondErr(w, err)
		return
	}
	respond(w, f, http.StatusCreated)
}

func filePatch(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID, fileID)
	if err != nil {
		respondErr(w, err)
		return
	}
	if err := assertJSON(r); err != nil {
		respondErr(w, err)
		return
	}
	data, err := ioutil.ReadAll(io.LimitReader(r.Body, defaultFileBodyLimit))
	if err != nil {
		respondErr(w, err)
		return
	}

	f, err := model.UpdateFile(r.Context(), ids[0], ids[1], ids[2], func(f *model.File) error {
		orig := *f
		if err := json.Unmarshal(data, f); err != nil {
			return err
		}
		f.ID = orig.ID
		f.Created = orig.Created
		f.Modified = time.Now().Truncate(time.Second)
		f.Project = orig.Project
		f.Commit = orig.Commit
		return nil
	})
	if err != nil {
		respondErr(w, err)
		return
	}
	respondOK(w, f)
}

func fileDelete(w http.ResponseWriter, r *http.Request) {
	ids, err := getIDs(r, projectID, commitID, fileID)
	if err != nil {
		respondErr(w, err)
		return
	}
	err = model.DeleteFile(r.Context(), ids[0], ids[1], ids[2])
	if err != nil {
		respondErr(w, err)
		return
	}
	respond(w, nil, http.StatusNoContent)
}
