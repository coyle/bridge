package files

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// List returns all file descriptors in a bucket
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// GetID retrieves the file ID from a bucket
func GetID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Get retrieves the file from a bucket
func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Delete removes the file from a bucket
func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// GetInfo retrieves the info for a file from a bucket
func GetInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// CreateEntryFromFrame __
func CreateEntryFromFrame(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// ListMirrorsForFile __
func ListMirrorsForFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
