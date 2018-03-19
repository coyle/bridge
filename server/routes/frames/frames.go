package frames

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Create initializes a new frame
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// AddShard adds an additional shard to a frame
func AddShard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// RemoveByID deletes a frame with the provided ID
func RemoveByID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Get retrieves all frames
func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// GetByID retrieves a frame with the provided ID
func GetByID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
