package health

import "github.com/emicklei/go-restful"

// Resource handles all health requests
type Resource struct{}

// ConfigureResource configures a new health resource
func ConfigureResource(container *restful.Container) *Resource {
	r := new(Resource)
	ws := new(restful.WebService)
	ws.
		Path("/health").
		Doc("Display Health Measurements").
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(r.displayMeasurements).
		Doc("display basic health").
		Operation("displayMeasurements").
		Writes(measurements{}))

	container.Add(ws)
	return r
}

func (h *Resource) displayMeasurements(request *restful.Request, response *restful.Response) {
	response.WriteEntity(takeMeasurements())
}
