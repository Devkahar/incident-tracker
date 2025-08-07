package router

import (
	"incident-tracker/config"
	"incident-tracker/controller"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	appContext *config.ApplicationContext
	router     *gin.Engine
}

func (server *Server) AddRoutes() *Server {
	v1 := server.router.Group("/v1")
	{
		incidents := v1.Group("/incidents")
		{
			incidents.POST("/", HandleRequest(server.appContext, controller.CreateIncident))
			incidents.GET("/", HandleRequest(server.appContext, controller.GetIncidents))
		}
	}
	return server
}

func NewServer(applicationContext *config.ApplicationContext) *Server {
	server := &Server{appContext: applicationContext, router: gin.Default()}
	server.router.Use(CORSMiddleware())
	return server
}

// Start starts the HTTP server.
func (s *Server) Start() *http.Server {
	server := &http.Server{
		Addr:              ":" + s.appContext.Config.Server.Port,
		Handler:           s.router,
		ReadHeaderTimeout: time.Second * 200,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("unable to start the server: ", err)
		}
	}()

	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
