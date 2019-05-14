package handle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VolticFroogo/Launcher-Server/db"
	"github.com/gorilla/mux"
)

type response struct {
	ID                            int
	Name, Path, Version, Download string   `json:",omitempty"`
	Exceptions                    []string `json:",omitempty"`
}

const (
	Update = iota
	Latest
)

// Listen to all incoming HTTP requests.
func Listen() {
	// Create a new router.
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Handle("/api/program", http.HandlerFunc(APIGetProgram)).Methods(http.MethodGet)

	// Start listening on the provided port.
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

// APIGetProgram handles all API requests for information on whether a program needs updating or not.
func APIGetProgram(w http.ResponseWriter, r *http.Request) {
	idParam, ok := r.URL.Query()["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var version string
	initial := false

	versionParam, ok := r.URL.Query()["version"]
	if !ok {
		_, ok := r.URL.Query()["initial"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		initial = true
	} else {
		version = versionParam[0]
	}

	id := idParam[0]

	program, err := db.GetProgram(id)
	if err != nil {
		log.Printf("Error getting program from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	latest := program.Latest()

	if !initial && version == latest {
		Respond(w, response{
			ID: Latest,
		})

		return
	}

	download := fmt.Sprintf("https://storage.googleapis.com/froogo-launcher/programs/%v/%v.zip", program.ID, latest)

	Respond(w, response{
		ID:         Update,
		Name:       program.Name,
		Path:       program.Path,
		Version:    latest,
		Download:   download,
		Exceptions: program.Exceptions,
	})
}

// Respond will JSON marshal data provided and write it to the stream.
func Respond(w http.ResponseWriter, data interface{}) {
	// Marshal the data into JSON bytes..
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	// Write the bytes to the stream.
	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Error writing bytes to stream: %v", err)
	}
}
