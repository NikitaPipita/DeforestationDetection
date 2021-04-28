package apiserver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func (s *server) MakeAndDownloadDump(w http.ResponseWriter, r *http.Request) {
	dumpRecordFilePath := s.store.Dump().CreateDump()
	if _, err := os.Stat(dumpRecordFilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(dumpRecordFilePath)))
	http.ServeFile(w, r, dumpRecordFilePath)
	_ = os.Remove(dumpRecordFilePath)
}

func (s *server) UploadAndExecuteDump(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	defer func() { _ = file.Close() }()
	qBytes, err := ioutil.ReadAll(file)
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	if err := s.store.Dump().Execute(string(qBytes)); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.respond(w, r, http.StatusOK, nil)
}
