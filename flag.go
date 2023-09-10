package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type FlagService struct {
	flags map[string]map[string]int
	mu    sync.Mutex
}

type Flag struct {
	Key               string `json:"key"`
	OnRolloutPercent  int    `json:"OnRolloutPercent"`
	OffRolloutPercent int    `json:"OffRolloutPercent"`
}

func NewFlagService() *FlagService {
	return &FlagService{
		flags: make(map[string]map[string]int),
	}
}

func (fs *FlagService) CreateFlag(w http.ResponseWriter, r *http.Request) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	flag := Flag{}

	err := json.NewDecoder(r.Body).Decode(&flag)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, exists := fs.flags[flag.Key]; exists {
		http.Error(w, "Flag already exists", http.StatusConflict)
		return
	}

	fs.flags[flag.Key] = map[string]int{
		"on":  flag.OnRolloutPercent,
		"off": flag.OffRolloutPercent,
	}

	w.WriteHeader(http.StatusCreated)
}

func (fs *FlagService) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	flag := Flag{}

	err := json.NewDecoder(r.Body).Decode(&flag)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, exists := fs.flags[flag.Key]; !exists {
		http.Error(w, "Flag does not exist", http.StatusNotFound)
		return
	}

	fs.flags[flag.Key] = map[string]int{
		"on":  flag.OnRolloutPercent,
		"off": flag.OffRolloutPercent,
	}

	w.WriteHeader(http.StatusOK)
}

func (fs *FlagService) GetUserFlagValue(w http.ResponseWriter, r *http.Request) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	userID := r.URL.Query().Get("user_id")
	flagKey := r.URL.Query().Get("flag_key")

	if flag, exists := fs.flags[flagKey]; exists {
		userIDInt := 0
		_, err := fmt.Sscanf(userID, "%d", &userIDInt)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		randomNumber := userIDInt % 100

		OnRolloutPercent := flag["on"]
		if randomNumber < OnRolloutPercent {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("on")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("off")
	} else {
		http.Error(w, "Flag not found", http.StatusNotFound)
	}

	log.Println(fs.flags)
}

func (fs *FlagService) ListFlags(w http.ResponseWriter, r *http.Request) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fs.flags)
}
