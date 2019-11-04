package main

import (
	"github.com/geraev/gokvserver/httpserver"
	"github.com/geraev/gokvserver/mapbased"
)

func main() {

	accounts := map[string]string{
		"iqoption": "qwerty64",
		"geraev":   "markus14",
	}

	ttt := mapbased.TestTestStorage()
	cache := mapbased.NewStorage()
	_ = cache

	http := httpserver.NewServer("8081", accounts, ttt)
	http.Run()

}

//TODO Заменить типы string, []string, map[string]string  на собственные алиасы этих типов
//TODO Вынести error-ы сервиса в errors.go
