package api

import (
	"event-data-pipeline/pkg/cli"

	"github.com/gin-gonic/gin"
)

type Service struct {
	port     int
	address  string
	scheme   string
	basePath string
	router   *gin.Engine
}

func NewService() *Service {

	mode := "release"
	if cli.Args.DebugEnabled {
		mode = "debug"
	}

	gin.SetMode(mode)
	// Create a gin router
	router := gin.New()

	svc := &Service{
		port:     cli.Args.Port,
		address:  cli.Args.Addr,
		scheme:   cli.Args.Scheme,
		basePath: cli.Args.BasePath,
		router:   router,
	}

	NewRouteHandler(svc)

	return svc
}
