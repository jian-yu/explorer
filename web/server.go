package web

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Daemon interface {
	Start() error
	Stop()
}

type server struct {
	addr string
	*gin.Engine
}

func NewServer(engine *gin.Engine) Daemon {
	addr := viper.GetString(`Web.Addr`)
	if addr == "" {
		addr = ":8080"
	}

	return &server{
		Engine: engine,
		addr:   addr,
	}
}

func (s *server) Start() error {
	return s.Run(s.addr)
}

func (s *server) Stop(){

}
