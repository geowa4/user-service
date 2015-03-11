package user

import (
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

// Resource handles all user requests
type Resource struct {
	dao map[string]user
}

// ConfigureResource configures a new health resource
func ConfigureResource(container *restful.Container) *Resource {
	u := new(Resource)
	u.dao = make(map[string]user)
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Doc("Display User Data").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(u.listUsers).
		Doc("list all users").
		Operation("listUsers").
		Writes([]user{}))

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(user{}))

	ws.Route(ws.POST("").To(u.createUser).
		Doc("create a user").
		Operation("createUser").
		Reads(user{}))

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		Doc("delete a user").
		Operation("removeUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	container.Add(ws)
	return u
}

func (u *Resource) listUsers(request *restful.Request, response *restful.Response) {
	numUsers := len(u.dao)
	if numUsers == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Users could not be found.")
		return
	}
	users := make([]user, 0, numUsers)
	for _, v := range u.dao {
		users = append(users, v)
	}
	response.WriteEntity(users)
}

func (u *Resource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	user, ok := u.dao[id]
	if !ok {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
		return
	}
	response.WriteEntity(user)
}

func (u *Resource) createUser(request *restful.Request, response *restful.Response) {
	user := new(user)
	err := request.ReadEntity(user)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	user.ID = strconv.Itoa(len(u.dao) + 1) // simple id generation
	u.dao[user.ID] = *user
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(user)
}

func (u *Resource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	user, ok := u.dao[id]
	if !ok {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
		return
	}
	delete(u.dao, id)
	response.WriteHeader(http.StatusOK)
	response.WriteEntity(user)
}
