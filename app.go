package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/redis.v5"
)

type App struct {
	Router      *mux.Router
	RedisClient *redis.Client
}

func (a *App) Initialize(port string, password string, dbStr string) {
	if len(port) == 0 {
		port = "6379"
	}

	if len(dbStr) == 0 {
		dbStr = "0"
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Fatal(err)
	}

	a.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", port),
		Password: password,
		DB:       db,
	})

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Printf("Listening at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/calendar/{building:[a-zA-Z]+}/{room:[0-9a-zA-Z]+}", a.getRoom).Methods("GET")
	a.Router.HandleFunc("/calendar/{building:[a-zA-Z]+}", a.getBuilding).Methods("GET")
}

func (a *App) getBuilding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	building, err := getBuilding(a.RedisClient, strings.ToUpper(vars["building"]))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, building)
}

func (a *App) getRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	room, err := getRoom(a.RedisClient, strings.ToUpper(vars["building"]), vars["room"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, room)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
