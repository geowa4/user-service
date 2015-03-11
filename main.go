package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/geowa4/user-service/health"
	"github.com/geowa4/user-service/user"
)

func main() {
	container := restful.NewContainer()
	health.ConfigureResource(container)
	user.ConfigureResource(container)

	config := swagger.Config{
		WebServices:     container.RegisteredWebServices(),
		WebServicesUrl:  "http://localhost:8080",
		ApiPath:         "/",
		SwaggerPath:     "/docs/",
		SwaggerFilePath: "/opt/swagger-ui",
	}
	swagger.RegisterSwaggerService(config, container)

	log.Printf("start listening on port 8080")
	server := &http.Server{Addr: ":8080", Handler: container}
	log.Fatal(server.ListenAndServe())
}
