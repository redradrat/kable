package api

import (
	"github.com/labstack/echo/v4"
)

type Serv struct{}

func StartUp(bind string) {
	serv := Serv{}
	e := echo.New()
	RegisterHandlers(e, &serv)
	e.Static("/", "kable.v1.yaml")
	e.Logger.Fatal(e.Start(bind))
}

func RegisterHandlers(e *echo.Echo, serv *Serv) {
	e.GET("/concepts", serv.GetConcepts)
	e.GET("/concept/:id", serv.GetConcept)
	e.GET("/repositories", serv.GetRepositories)
	e.GET("/repository/:id", serv.GetRepository)
	e.PUT("/repository/:id", serv.PutRepository)
	e.GET("/repository/:id/concepts", serv.GetRepositoryConcepts)
	e.GET("/repository/:id/concept/:path", serv.GetRepositoryConcept)
}
