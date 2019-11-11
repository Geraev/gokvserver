package httpserver

import (
	"errors"
	"github.com/geraev/gokvserver/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SetStringBody struct {
	Value string `json:"value" binding:"required"`
}

type SetTTLBody struct {
	Value uint64 `json:"value" binding:"required"`
}

type SetListBody struct {
	Value []string `json:"value" binding:"required"`
}

type SetDictionaryBody struct {
	Value map[string]string `json:"value" binding:"required"`
}

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

func (s *Server) Run() error {

	r := gin.Default()

	// Базовая аутентификация. Можно заменить на OAuth
	authorized := r.Group("/cache", gin.BasicAuth(s.accounts))

	authorized.GET("/keys", s.getKeys)
	authorized.GET("/key/:key", s.getElement)
	authorized.GET("/key/:key/:internalKey", s.getInternalElement)

	authorized.POST("/set/ttl/:key", s.setTTL)

	authorized.PUT("/set/string/:key", s.setString)
	authorized.PUT("/set/list/:key", s.setList)
	authorized.PUT("/set/dictionary/:key", s.setDictionary)

	authorized.DELETE("/remove/:key", s.deleteKey)

	return r.Run(":" + s.port)
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
			gin.H{"error": err.Error()}, //TODO Нельзя возвращать внутренние ошибки. Исправить в следующей реализации
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

// setTTL установка времени жизни ключа
// curl -H 'content-type: application/json' -k -u user:pass -d '{ "value": 3000 }' -X PUT http://localhost:8081/cache/set/ttl/<key>
func (s *Server) setTTL(c *gin.Context) {
	key := c.Param("key")
	var value SetTTLBody
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	s.storage.SetTTL(key, value.Value)
	return
}

// setSting добавление или обновление ключа строки в кеше
// curl -H 'content-type: application/json' -k -u user:pass -d '{ "value": "manu" }' -X PUT http://localhost:8081/cache/set/string/<key>
func (s *Server) setString(c *gin.Context) {
	key := c.Param("key")
	var value SetStringBody
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	s.storage.PutOrUpdateString(key, value.Value)
}

// setList добавление или обновление ключа списка в кеше
// curl -H 'content-type: application/json' -k -u user:pass -d '{ "value": ["manu","suro","jonk"] }' -X PUT http://localhost:8081/cache/set/list/<key>
func (s *Server) setList(c *gin.Context) {
	key := c.Param("key")
	var value SetListBody
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	s.storage.PutOrUpdateList(key, value.Value)
}

// setDictionary добавление или обновление ключа словаря в кеше
// curl -H 'content-type: application/json' -k -u user:pass -d '{"value": {"k1":"manu","k2":"sol","k3":"vano"} }' -X PUT http://localhost:8081/cache/set/dictionary/<key>
func (s *Server) setDictionary(c *gin.Context) {
	key := c.Param("key")
	var value SetDictionaryBody
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	s.storage.PutOrUpdateDictionary(key, value.Value)
}

// deleteKey удаление ключа из кеша
// curl -k -u user:pass -X DELETE http://localhost:8081/cache/remove/<key>
func (s *Server) deleteKey(c *gin.Context) {
	key := c.Param("key")
	s.storage.RemoveElement(key)
}
