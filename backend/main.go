package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    api := r.PathPrefix("/v1").Subrouter()

    api.HandleFunc("/auth/login", loginHandler).Methods("POST")
    api.HandleFunc("/services", servicesHandler).Methods("GET")
    api.HandleFunc("/instances", createInstanceHandler).Methods("POST")
    api.HandleFunc("/instances/{id}/start", startInstanceHandler).Methods("POST")
    api.HandleFunc("/storage/kv/{key}", kvGetHandler).Methods("GET")

    // WebSocket endpoint for IPC
    r.HandleFunc("/v1/sessions/{sid}/ws", wsHandler)

    // serve frontend dist (relative)
    dist := filepath.Join("..", "frontend", "dist")
    if _, err := os.Stat(dist); os.IsNotExist(err) {
        log.Println("Warning: frontend dist missing, run frontend build or run dev server.")
    }
    r.PathPrefix("/").Handler(http.FileServer(http.Dir(dist)))

    port := getEnv("PORT", "8080")
    log.Printf("Listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(k, fallback string) string {
    if v := os.Getenv(k); v != "" { return v }
    return fallback
}