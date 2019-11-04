package main

import (
	"github.com/geraev/gokvserver/httpserver"
	"github.com/geraev/gokvserver/mapbased"
)

func main() {
	ttt := mapbased.NewStorage()
	_ = ttt


	accounts := map[string]string{
		"iqoption": "qwerty64",
		"geraev":   "markus14",
	}

	http := httpserver.NewServer("8081", accounts)
	http.Run()

}


