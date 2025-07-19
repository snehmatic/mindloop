package v1

// define api logic here
import (
	"net/http"

	"github.com/snehmatic/mindloop/internal/utils"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse([]byte("Welcome to the Mindloop!"), w, http.StatusOK)
}

func HandleHealthz(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse([]byte("OK"), w, http.StatusOK)
}
