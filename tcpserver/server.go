package tcpserver

import (
	"encoding/json"
	"fmt"
	"github.com/bsm/redeo"
	"github.com/bsm/redeo/resp"
	"github.com/geraev/gokvserver/structs"
	"log"
	"net"
)

const (
	errSetMsg = `Set or update value
Examples:
  set string new_key string_value
  set list planets '{"value": ["earth","jupiter","saturn"], "ttl": 10000}'
  set dictionary planets_map '{"value": ["earth":2220,"jupiter":3899,"saturn":23000], "ttl": 10000}'
`
)

type BodyList struct {
	Value []string `json:"value" binding:"required"`
}

type BodyDictionary struct {
	Value map[string]string `json:"value" binding:"required"`
}

type SetTTLBody struct {
	Value uint64 `json:"value" binding:"required"`
}

type Server struct {
	port    string
	storage structs.Storage
}

func NewServer(port string, storage structs.Storage) *Server {
	return &Server{
		port:    port,
		storage: storage,
	}
}

func (s *Server) Run() error {
	srv := redeo.NewServer(nil)
	srv.Handle("ping", redeo.Ping())
	srv.Handle("echo", redeo.Echo())
	srv.Handle("info", redeo.Info(srv))

	srv.HandleFunc("set", s.set)

	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	defer lis.Close()

	log.Printf("waiting for connections on %s", lis.Addr().String())
	return srv.Serve(lis)
}

// getKeys получение списка ключей из кеша
// curl -k -u user:pass http://localhost:8081/cache/keys
func (s *Server) getKeys(w resp.ResponseWriter, c *resp.Command) {
	/*	c.JSON(
			http.StatusOK,
			gin.H{"keys": s.storage.GetKeys()},
		)
	*/
}

// getKeys получение элемента из кеша по ключу
// curl -k -u user:pass http://localhost:8081/cache/key/<key>
func (s *Server) getElement(w resp.ResponseWriter, c *resp.Command) {
	/*	key := c.Param("key")

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
	*/
}

// getInternalElement получение внутреннего элемента из списка (по индексу) или словаря (по ключу)
// curl -k -u user:pass http://localhost:8081/cache/key/<key>/<internal key or index>
func (s *Server) getInternalElement(w resp.ResponseWriter, c *resp.Command) {
	/*	key := c.Param("key")
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
	*/
}

// setTTL установка времени жизни ключа
// curl -H 'content-type: application/json' -k -u user:pass -d '{ "value": 3000 }' -X PUT http://localhost:8081/cache/set/ttl/<key>
func (s *Server) setTTL(w resp.ResponseWriter, c *resp.Command) {
	/*	key := c.Param("key")
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
	*/
}

// set добавление или обновление ключа в кеше
// curl -H 'content-type: application/json' -k -u user:pass -d '{ "value": "manu", "ttl": 5000 }' -X PUT http://localhost:8081/cache/set/string/<key>
func (s *Server) set(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() < 3 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		w.AppendInt(int64(c.ArgN()))
		w.AppendError(errSetMsg)
		return
	}

	var (
		vartype   = c.Arg(0).String()
		key       = c.Arg(1).String()
		val       []byte
		isUpdated bool
	)

	for _, item := range c.Args[2:] {
		val = append(append(val, " "...), item.Bytes()...)
	}

	switch vartype {
	case "string":
		_, isUpdated = s.storage.PutOrUpdateString(key, string(val))
	case "list":
		var value BodyList
		err := json.Unmarshal(val, &value)
		if err != nil {
			w.AppendError(err.Error())
			return
		}
		_, isUpdated = s.storage.PutOrUpdateList(key, value.Value)
	case "dictionary":
		var value BodyDictionary
		err := json.Unmarshal(val, &value)
		if err != nil {
			w.AppendError(err.Error())
			return
		}
		_, isUpdated = s.storage.PutOrUpdateDictionary(key, value.Value)
	default:
		w.AppendError(redeo.UnknownCommand(c.Name))
		w.AppendError(errSetMsg)
		return
	}

	if isUpdated {
		w.AppendError(fmt.Sprintf("key %s was updated", key))
		return
	} else {
		w.AppendInlineString(fmt.Sprintf("key %s was set", key))
	}

	w.AppendOK()
}

/*func (s *Server) setSting(w resp.ResponseWriter, c *resp.Command) {
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
	if value.TTL != 0 {
		s.storage.SetTTL(key, value.TTL)
	}
}
*/
// setList добавление или обновление ключа списка в кеше
// curl -H 'content-type: application/json' -k -u user:pass -d '{ "value": ["manu","suro","jonk"] }' -X PUT http://localhost:8081/cache/set/list/<key>
func (s *Server) setList(w resp.ResponseWriter, c *resp.Command) {
	/*	key := c.Param("key")
		var value SetListBody
		if err := c.ShouldBindJSON(&value); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
		s.storage.PutOrUpdateList(key, value.Value)
	*/
}

// setDictionary добавление или обновление ключа словаря в кеше
// curl -H 'content-type: application/json' -k -u user:pass -d '{"value": {"k1":"manu","k2":"sol","k3":"vano"} }' -X PUT http://localhost:8081/cache/set/dictionary/<key>
func (s *Server) setDictionary(w resp.ResponseWriter, c *resp.Command) {
	/*	key := c.Param("key")
		var value SetDictionaryBody
		if err := c.ShouldBindJSON(&value); err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
		s.storage.PutOrUpdateDictionary(key, value.Value)
	*/
}

// deleteKey удаление ключа из кеша
// curl -k -u user:pass -X DELETE http://localhost:8081/cache/remove/<key>
func (s *Server) deleteKey(w resp.ResponseWriter, c *resp.Command) {
	/*	key := c.Param("key")
		s.storage.RemoveElement(key)
	*/
}
