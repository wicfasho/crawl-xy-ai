package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wicfasho/crawl-xy/sqlc"
)

type Server interface {
	GetDBStore() *sqlc.Queries
	GetRouter() *gin.Engine
}

var server Server

func SetServer(s Server) {
	server = s
}

func GetServer() Server {
	return server
}
