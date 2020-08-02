package web

import (
	"os"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	env := "GIN_LOGGER"
	engine := gin.New()
	engine.Use(gin.Recovery())

	if os.Getenv(env) == "1"{
		engine.Use(gin.Logger())
	}
	return engine
}
