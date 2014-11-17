// Copyright 2014 Bowery, Inc.
// Contains the routes for mercer.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Bowery/gopackages/requests"
	"github.com/Bowery/gopackages/sys"
	"github.com/Bowery/gopackages/tar"
	"github.com/Bowery/gopackages/web"
	"github.com/unrolled/render"
)

var Routes = []web.Route{
	{"GET", "/", HelloHandler, false},
	{"GET", "/healthz", HealthzHandler, false},
	{"POST", "/code", AnalyzeCodeHandler, false},
}

var renderer = render.New(render.Options{
	IndentJSON:    true,
	IsDevelopment: true,
})

type Commands struct {
	Start string `json:"start"`
	Build string `json:"build"`
	Init  string `json:"init"`
}

var LanguageToCommands = map[string]*Commands{
	"JavaScript": &Commands{
		Start: "npm start",
		Build: "npm install",
		Init:  "apt-get install -y nodejs",
	},
}

func HelloHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "Bowery Code Analyzer")
}

func HealthzHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "ok")
}

func AnalyzeCodeHandler(rw http.ResponseWriter, req *http.Request) {
	tarball, _, err := req.FormFile("file")
	if err != nil {
		renderer.JSON(rw, http.StatusBadRequest, map[string]string{
			"status": requests.STATUS_FAILED,
			"error":  err.Error(),
		})
		return
	}
	defer tarball.Close()

	analysisPath := filepath.Join(os.Getenv(sys.HomeVar), ".mercer", uuid.New())

	if err = os.MkdirAll(analysisPath, os.ModePerm|os.ModeDir); err != nil {
		renderer.JSON(rw, http.StatusInternalServerError, map[string]string{
			"status": requests.STATUS_FAILED,
			"error":  err.Error(),
		})
		return
	}

	if err := tar.Untar(tarball, analysisPath); err != nil {
		renderer.JSON(rw, http.StatusInternalServerError, map[string]string{
			"status": requests.STATUS_FAILED,
			"error":  err.Error(),
		})
		return
	}

	linguistPath, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "classifiers/language"))
	linguistOut, err := exec.Command(linguistPath, analysisPath).Output()
	if err != nil {
		renderer.JSON(rw, http.StatusInternalServerError, map[string]string{
			"status": requests.STATUS_FAILED,
			"error":  err.Error(),
		})
		return
	}
	languages := map[string]interface{}{}
	if err := json.Unmarshal(linguistOut, &languages); err != nil {
		renderer.JSON(rw, http.StatusInternalServerError, map[string]string{
			"status": requests.STATUS_FAILED,
			"error":  err.Error(),
		})
		return
	}

	cmds := &Commands{}
	for language, weight := range languages {
		fmt.Println(language, "-", weight)
		lc := LanguageToCommands[language]
		if lc != nil {
			cmds.Start += lc.Start + "\n"
			cmds.Build += lc.Build + "\n"
			cmds.Init += lc.Init + "\n"
		} else {
			fmt.Println(language, "is not currently supported")
		}

	}

	renderer.JSON(rw, http.StatusOK, map[string]interface{}{
		"status":   requests.STATUS_SUCCESS,
		"commands": cmds,
	})
}
