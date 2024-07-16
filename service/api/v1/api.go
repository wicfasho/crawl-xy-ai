package api

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wicfasho/crawl-xy/routes"
	"github.com/wicfasho/crawl-xy/server"
	"github.com/wicfasho/crawl-xy/sqlc"
)

type Server struct {
	dbStore *sqlc.Queries
	router  *gin.Engine
}

var apiVersion = "v1"

func NewServer(dbStore *sqlc.Queries) (Server, error) {
	s := Server{
		dbStore: dbStore,
	}

	s.setupRouter()

	server.SetServer(&s)

	return s, nil
}

func (server *Server) setupRouter() error {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(r.Group("api/" + apiVersion))

	server.router = r

	return nil
}

func (server *Server) GetDBStore() *sqlc.Queries {
	return server.dbStore
}

func (server *Server) GetRouter() *gin.Engine {
	return server.router
}

func (server *Server) Start(wg *sync.WaitGroup) {
	wg.Wait()
	server.router.Run()
}
