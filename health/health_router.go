package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/geowa4/user-service/wrappers"
	"github.com/gorilla/mux"
)

const routerName = "health"

// GetName returns the name for this router.
func GetName() string {
	return routerName
}

// Router encapsulates the state required for handling user requests.
type Router struct {
	Subrouter *mux.Router
	DB        *sql.DB
}

// NewRouter makes a new Router with its DB set to the argument.
func NewRouter(db *sql.DB, sub *mux.Router) (hr *Router) {
	hr = new(Router)
	hr.Subrouter = sub
	hr.DB = db
	return
}

// HandleRoutes configures the provided router to handle health requests
func (hr *Router) HandleRoutes() {
	hr.Subrouter.
		Methods("GET").
		Path("/").
		Name("SimpleHealth").
		Handler(wrappers.Defaults(hr.healthHandler, "SimpleHealth"))
}

// Measurements of the system.
type Measurements struct {
	NumCPU          int    `json:"num_cpu"`
	NumGoroutine    int    `json:"num_goroutine"`
	MemoryAllocated uint64 `json:"memory_allocated"`
	DBConnected     bool   `json:"db_connected"`
	DBPingError     string `json:"db_ping_error,omitempty"`
}

// GetSimpleMeasurements returns some simple measurements about the system.
func (hr *Router) GetSimpleMeasurements() Measurements {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	// hr.DB.Ping() doesn't work so just select something arbitrary.
	_, pingErr := hr.DB.Exec("SELECT 1 AS ping")
	pingErrMsg := ""
	if pingErr != nil {
		pingErrMsg = pingErr.Error()
	}
	return Measurements{
		NumCPU:          runtime.NumCPU(),
		NumGoroutine:    runtime.NumGoroutine(),
		MemoryAllocated: memStats.Alloc,
		DBConnected:     pingErr == nil,
		DBPingError:     pingErrMsg,
	}
}

func (hr *Router) healthHandler(res http.ResponseWriter, req *http.Request) {
	measurements := hr.GetSimpleMeasurements()
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(res).Encode(measurements)
}
