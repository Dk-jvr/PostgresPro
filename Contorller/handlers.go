package Contorller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"postgres/Database"
	"postgres/Models"
	"strconv"
)

type (
	ScriptManager struct {
		Scripts map[int64]*exec.Cmd
	}
)

func (manager *ScriptManager) CreateScript(writer http.ResponseWriter, request *http.Request) {
	var scriptId int64
	data := new(Models.ScriptData)
	err := json.NewDecoder(request.Body).Decode(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()

	err = validate.Struct(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	scriptId, err = Database.AddScript(*data)
	data.Id = scriptId
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command(data.Type, "-c", data.Script)
	manager.Scripts[scriptId] = cmd
	errCh := make(chan error, 1)
	go func() {
		var buffer bytes.Buffer
		cmd.Stdout = &buffer
		err = cmd.Start()
		if err != nil {
			log.Println(err)
			errCh <- err
			return
		}
		errCh <- err
		if err = cmd.Wait(); err != nil {
			buffer.WriteString(fmt.Sprintf("%s: command was stopped", err))
		}
		err = Database.InsertScriptData(*data, buffer.String())
		if err != nil {
			log.Fatal(err)
		}
		delete(manager.Scripts, scriptId)
		return
	}()
	select {
	case err = <-errCh:
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
	return
}

func GetScript(writer http.ResponseWriter, request *http.Request) {
	var response []byte
	scriptData := new(Models.ScriptData)
	vars := mux.Vars(request)
	scriptId, err := strconv.ParseInt(vars["script_id"], 10, 64)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, fmt.Sprintf("Something wrong: %v", err), http.StatusInternalServerError)
		} else {
			http.Error(writer, "No scripts found", http.StatusNotFound)
		}
		return
	}
	*scriptData, err = Database.GetScript(scriptId)
	response, err = json.Marshal(scriptData)
	writer.Write(response)
}

func GetScripts(writer http.ResponseWriter, request *http.Request) {
	var (
		scripts  []Models.ScriptData
		response []byte
		err      error
	)
	scripts, err = Database.GetScripts()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			response, _ = json.Marshal(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			writer.Write(response)
		} else {
			http.Error(writer, "No scripts found", http.StatusNotFound)
		}
		return
	}
	response, err = json.Marshal(scripts)
	writer.Write(response)
	return
}

func (manager *ScriptManager) ExecuteScript(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	var (
		script Models.ScriptData
	)
	scriptId, err := strconv.ParseInt(vars["script_id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	script, err = Database.GetScript(scriptId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, "Script not found", http.StatusNotFound)
		}
		return
	}
	cmd := exec.Command(script.Type, "-c", script.Script)
	manager.Scripts[scriptId] = cmd
	errCh := make(chan error, 1)
	go func() {
		var buffer bytes.Buffer
		cmd.Stdout = &buffer
		err = cmd.Start()
		if err != nil {
			log.Println(err)
			errCh <- err
			return
		}
		errCh <- err
		if err = cmd.Wait(); err != nil {
			buffer.WriteString(fmt.Sprintf("%s: command was stopped", err))
		}
		err = Database.InsertScriptData(script, buffer.String())
		if err != nil {
			log.Fatal(err)
		}
		delete(manager.Scripts, scriptId)
		return
	}()
	select {
	case err = <-errCh:
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
	return
}

func (manager *ScriptManager) TerminateScript(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	scriptId, err := strconv.ParseInt(vars["script_id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	if _, ok := manager.Scripts[scriptId]; ok {
		manager.Scripts[scriptId].Process.Kill()
		delete(manager.Scripts, scriptId)
	} else {
		http.Error(writer, "Active script not found", http.StatusNotFound)
	}
}

func UserAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("token")
		if token != "user_token" && token != "admin_token" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		if token == "user_token" && (request.URL.Path != "/commands/select" || request.URL.Path != "/commands/select/{command_id}") {
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
