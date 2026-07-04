package main

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/swaggo/gin-swagger/example/multiple/api/v1"
	v2 "github.com/swaggo/gin-swagger/example/multiple/api/v2"
	_ "github.com/swaggo/gin-swagger/example/multiple/docs"
	ginScalar "github.com/tinyauthapp/gin-scalar"
)

func main() {
	// New gin router
	router := gin.New()

	// Register api/v1 endpoints
	v1.Register(router)
	router.GET("/scalar/v1/*any", ginScalar.WrapHandler(nil, ginScalar.InstanceName("v1"), ginScalar.BasePath("/scalar/v1")))

	// Register api/v2 endpoints
	v2.Register(router)
	router.GET("/scalar/v2/*any", ginScalar.WrapHandler(nil, ginScalar.InstanceName("v2"), ginScalar.BasePath("/scalar/v2")))

	// Listen and Server in
	_ = router.Run()
}
