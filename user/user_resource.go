package user

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

// Resource handles all user requests
type Resource struct {
	dao map[string]user
}

// ConfigureResource configures a new health resource
func ConfigureResource(container *restful.Container) *Resource {
	r := new(Resource)
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Doc("Display User Data").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(r.listUsers).
		Doc("list all users").
		Operation("listUsers").
		Writes([]user{}))

	container.Add(ws)
	return r
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
