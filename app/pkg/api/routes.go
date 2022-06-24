package api

import "github.com/gin-gonic/gin"

type Routes struct {
	router *gin.Engine
}

func NewRouteHandler(service *Service) *Routes {
	router := &Routes{
		router: service.router,
	}
	// router.setHealthGetHandler(service.GetConfig().BasePath)

	// if service.GetConfig().DebugEnabled {
	// 	router.setSwaggerGetHandler(
	// 		service.GetConfig().Scheme,
	// 		service.GetConfig().Addr,
	// 		service.GetConfig().Port,
	// 		service.GetConfig().BasePath)
	// }

	return router
}
