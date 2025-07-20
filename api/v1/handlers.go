package v1

// define api logic here
import (
	"net/http"

	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/utils"
)

type MindloopHandler struct {
	config *config.Config
	// db client todo
}

func NewMindloopHandler() *MindloopHandler {
	return &MindloopHandler{
		config: config.GetConfig(),
	}
}

func (mlh *MindloopHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse([]byte("Welcome to the Mindloop!"), w, http.StatusOK)
}

func (mlh *MindloopHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse([]byte("OK"), w, http.StatusOK)
}
