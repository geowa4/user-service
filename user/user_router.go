package user

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/geowa4/user-service/wrappers"
	"github.com/geowa4/user-service/user/github"
	"github.com/gorilla/mux"
)

const routerName = "users"

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
func NewRouter(db *sql.DB, sub *mux.Router) (ur *Router) {
	ur = new(Router)
	ur.Subrouter = sub
	ur.DB = db
	return
}

// HandleRoutes configures the provided router to handle user requests
func (ur *Router) HandleRoutes() {
	ur.Subrouter.
		Methods("POST").
		Path("/").
		Name("CreateUser").
		Handler(wrappers.Defaults(ur.createUser, "CreateUser"))
	ur.Subrouter.
		Methods("GET").
		Path("/").
		Name("ListUsers").
		Handler(wrappers.Defaults(ur.listUsers, "ListUsers"))
}

// ListResponse convenience struct for encapsulating list responses.
type ListResponse struct {
	Users Users `json:"users"`
}

func (ur *Router) listUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	users, err := LoadAll(ur.DB)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(ListResponse{
			Users: users,
		})
	}
}

// CreateRequest encapsulates the data sent by the client to create a user.
type CreateRequest struct {
	Code string `json:"code"`
}

func (ur *Router) createUser(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(res).Encode(err)
	} else {
		createReq := new(CreateRequest)
		json.NewDecoder(req.Body).Decode(createReq)
		ghUser, err := github.GetUserFromCode(createReq.Code)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(res).Encode(err)
		}
		user := User{
			Email:           ghUser.Email,
			Name:            ghUser.Name,
			AccessToken:     ghUser.AccessToken,
			Scope:           ghUser.Scope,
			GitHubID:        ghUser.ID,
			GitHubLogin:     ghUser.Login,
			GitHubAvatarURL: ghUser.AvatarURL,
			GitHubHTMLURL:   ghUser.HTMLURL,
		}
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := user.Save(ur.DB); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(err)
		} else {
			res.WriteHeader(http.StatusOK)
			json.NewEncoder(res).Encode(user)
		}
	}
}
