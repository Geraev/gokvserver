package tcpserver

import (
	"encoding/json"
	"fmt"
	"github.com/bsm/redeo"
	"github.com/bsm/redeo/resp"
	"github.com/geraev/gokvserver/structs"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	errSetMsg = `Set or update value
Examples:
  set string new_key string_value
  set list planets '{"value": ["earth","jupiter","saturn"], "ttl": 10000}'
  set dictionary planets_map '{"value": ["earth":2220,"jupiter":3899,"saturn":23000], "ttl": 10000}'
`

	errKeysMsg = `Get all keys
Example:
  keys
`

	errKeyMsg = `Get value for key (and internal key)
Example:
  key <key>
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
	srv.HandleFunc("keys", s.getKeys)
	srv.HandleFunc("key", s.getElement)
	srv.HandleFunc("ikey", s.getInternalElement)

	srv.HandleFunc("expire", s.expire)
	srv.HandleFunc("remove", s.deleteKey)

	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	defer lis.Close()

	log.Printf("waiting for connections on %s", lis.Addr().String())
	return srv.Serve(lis)
}

// getKeys получение списка ключей из кеша
func (s *Server) getKeys(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 0 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		w.AppendError(errKeysMsg)
		return
	}

	result := s.storage.GetKeys()
	w.AppendInlineString(strings.Join(result, ", "))
}

// getKey получение элемента из кеша по ключу
func (s *Server) getElement(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 1 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		w.AppendError(errKeyMsg)
		return
	}
	var (
		key   = c.Arg(0).String()
		value []byte
	)

	val, err := s.storage.GetElement(key)
	if err != nil {
		w.AppendError(err.Error())
		return
	}

	switch val.(type) {
	case string, []string, map[string]string:
		value, err = json.Marshal(val)
		if err != nil {
			w.AppendError(err.Error())
			return
		}
	default:
		w.AppendError("something wrong: type error")
		return
	}
	w.AppendInline(value)
}

// getInternalElement получение внутреннего элемента из списка (по индексу) или словаря (по ключу)
func (s *Server) getInternalElement(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 2 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		w.AppendError(errSetMsg)
		return
	}
	var (
		val         string
		key         = c.Arg(0).String()
		internalKey = c.Arg(1).String()
	)

	vartype, err := s.storage.GetType(key)
	if err != nil {
		w.AppendError(err.Error())
		return
	}

	switch vartype {
	case structs.List:
		index, err := strconv.ParseUint(internalKey, 10, 0)
		if err != nil {
			w.AppendError(err.Error())
			return
		}
		val, err = s.storage.GetListElement(key, int(index))
		if err != nil {
			w.AppendError(err.Error())
			return
		}
	case structs.Dictionary:
		val, err = s.storage.GetDictionaryElement(key, internalKey)
		if err != nil {
			w.AppendError(err.Error())
			return
		}
	default:
		w.AppendError("something wrong: type error")
		return
	}

	w.AppendInlineString(val)
}

// expire установка времени жизни ключа
func (s *Server) expire(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 2 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		return
	}
	var (
		key = c.Arg(0).String()
	)
	val, err := c.Arg(1).Int()
	if err != nil {
		w.AppendError(err.Error())
		return
	}

	s.storage.SetTTL(key, uint64(val))

	w.AppendOK()
}

// set добавление или обновление ключа в кеше
func (s *Server) set(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() < 3 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
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
	} else {
		w.AppendInlineString(fmt.Sprintf("key %s was set", key))
	}
}

// deleteKey удаление ключа из кеша
func (s *Server) deleteKey(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 1 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		return
	}

	var (
		key = c.Arg(0).String()
	)

	s.storage.RemoveElement(key)
}
