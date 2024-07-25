package main

import (
	"github.com/blanc42/ecms/pkg/initializers"
	"github.com/blanc42/ecms/pkg/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	routes.SetupRouter(r)
	r.Run(":8080")
}
