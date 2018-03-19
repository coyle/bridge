package buckets

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Get retrieves a bucket from the database
func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// GetByID retrieves a bucket with the provided ID
func GetByID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Create initializes a new bucket
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// DestroyByID removes a bucket with the provided ID
func DestroyByID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// UpdateByID updates a bucket with the provided ID
func UpdateByID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// CreateToken initializes a new token for the bucket associated with the provided ID
func CreateToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
