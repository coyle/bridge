package contacts

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetList retrieves a contact list
func GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// GetByNodeID retrieves a contact by the node ID
func GetByNodeID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// PatchByNodeID updates the contacts on a Node
func PatchByNodeID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Create initalizes a new contact
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// CreateChallenge initializes a new challenge for the contact
func CreateChallenge(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
