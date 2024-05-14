package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"postgres/Contorller"
	"postgres/Database"
)

func main() {

	manager := &Contorller.ScriptManager{
		Scripts: make(map[int64]*exec.Cmd),
	}
	database := Database.ConnectDatabase()

	defer database.Close()

	router := mux.NewRouter()
	router.HandleFunc("/commands/add", manager.CreateScript)
	router.HandleFunc("/commands/select", Contorller.GetScripts)
	router.HandleFunc("/commands/select/{script_id}", Contorller.GetScript)
	router.HandleFunc("/commands/execute/{script_id}", manager.ExecuteScript)
	router.HandleFunc("/commands/terminate/{script_id}", manager.TerminateScript)

	userAuthHandler := Contorller.UserAuthMiddleWare(router)

	log.Fatal(http.ListenAndServe(":8080", userAuthHandler))
}
