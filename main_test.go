package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"postgres/Contorller"
	"postgres/Database"
	"postgres/Models"
	"testing"
)

var (
	manager = &Contorller.ScriptManager{
		Scripts: make(map[int64]*exec.Cmd),
	}
	db *sql.DB
)

func TestInsertScript(t *testing.T) {
	var (
		err      error
		request  *http.Request
		response *http.Response
		data     = Models.ScriptData{Script: `#!/bin/bash\n\n   echo \"Some data\"\n`, Type: `/bin/bash`}
	)
	testServer := httptest.NewServer(http.HandlerFunc(manager.CreateScript))
	jsonData, _ := json.Marshal(data)
	defer testServer.Close()
	db = Database.ConnectDatabase()
	url := testServer.URL + "/commands/add"
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("token", "admin_token")
	if err != nil {
		t.Error(err)
	}
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("POST command failed with status %d", response.StatusCode)
	}
}
