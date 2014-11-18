// Copyright 2014 Bowery, Inc.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Bowery/gopackages/ignores"
	"github.com/Bowery/gopackages/tar"
)

var testApp io.Reader

func init() {
	// Create tar'd app to do tests with
	path := filepath.Join("test", "app")
	fmt.Println("Test app", path)
	ignoreList, err := ignores.Get(path)
	if err != nil {
		panic(err)
	}

	testApp, err = tar.Tar(path, ignoreList)
	if err != nil {
		panic(err)
	}
}

type mercerRes struct {
	Status string `json:"status"`
	Err    string `json:"error"`
	commands
}

func newAnalyzeCodeRequest(url string) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "upload")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, testApp)
	if err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, &body)
	if req != nil {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	return req, err
}

func TestAnalyzeCodeHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(analyzeCodeHandler))
	defer server.Close()

	req, err := newAnalyzeCodeRequest(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("Failed to analyze code. StatusCode:", res.StatusCode)
	}
	mercerResponse := new(mercerRes)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(mercerResponse)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(mercerResponse)
}
