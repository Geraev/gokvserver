package httpserver

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port     string
	accounts gin.Accounts
}

func NewServer(port string, accounts map[string]string) *Server {
	return &Server{
		port:     port,
		accounts: accounts,
	}
}

func (s *Server) Run() {

	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(s.accounts))

	authorized.GET("/keys", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	log.Logger.Fatal(r.Run(":" + s.port))
}

func (s *Server) getKeys(c *gin.Context) {
}

func (s *Server) setSting(c *gin.Context) {
}

func (s *Server) setList(c *gin.Context) {
}

func (s *Server) setDictionary(c *gin.Context) {
}

