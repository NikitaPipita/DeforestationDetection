package apiserver

import (
	"deforestation.detection.com/server/internal/app/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func (a *server) ServeDumpDownloadRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// ЕСЛИ ПАПКА ДАМПА НЕ СУЩЕСТВУЕТ,
		// CreateDir (путь папки дампа)

		dumpRecord, err := a.store.Dump().Make("dumps/")
		if err != nil {
			a.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		if err != nil || !filesutil.Exist(dumpRecord.FilePath) {
			a.server.Respond(w, r, http.StatusNotFound, nil)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", dumpRecord.FileName))
		http.ServeFile(w, r, dumpRecord.DumpFilePath)
	}
}

func Exist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
Delete function deletes the file, specified in the argument and returns
the boolean value of operation.

If the file is deleted, true is returned.
If the file is not exist, true will be returned as well.

Otherwise, the function returns false
*/
func Delete(filePath string) bool {
	if err := os.Remove(filePath); err != nil {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return true
		}
		return false
	}
	return true
}

// ExtractFileName selects the file name from it's path
func ExtractFileName(filePath string) string {
	return filepath.Base(filePath)
}

func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

func ClearDir(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := os.RemoveAll(path.Join(dirPath, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
