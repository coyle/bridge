package keys

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Add a new public key
func Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Get retrieves all public keys
func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Remove deletes a public key with the provided ID
func Remove(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
