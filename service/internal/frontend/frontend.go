package frontend

import (
	"net/http"

	"github.com/jamesread/golure/pkg/dirs"
)

func findWebuiDir() string {
	directoriesToSearch := []string{
		"../frontend/",
		"/app/frontend/",
	}

	dir, err := dirs.GetFirstExistingDirectory("webui", directoriesToSearch)

	if err != nil {
		panic("Failed to find webui directory: " + err.Error())
	}

	return dir
}

func GetNewHandler() (http.Handler) {
	return http.FileServer(http.Dir(findWebuiDir()))
}
