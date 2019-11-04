package httpserver

import (
	"errors"
	"github.com/geraev/gokvserver/structs"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port     string
	accounts gin.Accounts
	storage  structs.Storage
}

//TODO Вынести таблицу аккаунтов из обьекта Server
func NewServer(port string, accounts map[string]string, storage structs.Storage) *Server {
	return &Server{
		port:     port,
		accounts: accounts,
		storage:  storage,
	}
}

func (s *Server) Run() {

	r := gin.Default()

	// Базовая аутентификация. Можно заменить на OAuth
	authorized := r.Group("/cache", gin.BasicAuth(s.accounts))

	authorized.GET("/keys", s.getKeys)
	authorized.GET("/key/:key", s.getElement)
	authorized.GET("/key/:key/:internalKey", s.getInternalElement)

	log.Fatal(r.Run(":" + s.port))
}

// getKeys получение списка ключей из кеша
// curl -k -u user:pass http://localhost:8081/cache/keys
func (s *Server) getKeys(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{"keys": s.storage.GetKeys()},
	)
}

// getKeys получение элемента из кеша по ключу
// curl -k -u user:pass http://localhost:8081/cache/key/<key>
func (s *Server) getElement(c *gin.Context) {
	key := c.Param("key")

	val, err := s.storage.GetElement(key)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},  //TODO Нельзя возвращать внутренние ошибки. Исправить в следующей реализации
		)
		return
	}

	switch v := val.(type) {
	case string, []string, map[string]string:
		c.JSON(
			http.StatusOK,
			gin.H{"value": v},
		)
	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("something wrong: type error")},
		)
	}
}

// getInternalElement получение внутреннего элемента из списка (по индексу) или словаря (по ключу)
// curl -k -u user:pass http://localhost:8081/cache/key/<key>/<internal key or index>
func (s *Server) getInternalElement(c *gin.Context) {
	key := c.Param("key")
	internalKey := c.Param("internalKey")

	vartype, err := s.storage.GetType(key)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	var val string
	switch vartype {
	case structs.List:
		index, err := strconv.ParseUint(internalKey, 10, 0)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
		val, err = s.storage.GetListElement(key, int(index))
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
	case structs.Dictionary:
		val, err = s.storage.GetDictionaryElement(key, internalKey)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("something wrong: type error").Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"value": val},
	)
	return
}

func (s *Server) setSting(c *gin.Context) {

}

func (s *Server) setList(c *gin.Context) {
}

func (s *Server) setDictionary(c *gin.Context) {
}
