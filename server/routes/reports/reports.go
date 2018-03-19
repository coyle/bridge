package reports

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Create a new exchange report
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
