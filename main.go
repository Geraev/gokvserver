package main

import (
	"flag"
	"github.com/geraev/gokvserver/httpserver"
	"github.com/geraev/gokvserver/mapbased"
	"github.com/geraev/gokvserver/structs"
	"github.com/geraev/gokvserver/tcpserver"
	"log"
)

var (
	flags struct {
		tcpAddr  string
		httpAddr string
	}

	cache structs.Storage

	accounts = map[string]string{
		"iqoption": "qwerty64",
		"geraev":   "markus14",
	}
)

func init() {
	flag.StringVar(&flags.tcpAddr, "tcp-port", "9736", "The TCP port to bind to")
	flag.StringVar(&flags.httpAddr, "http-port", "8081", "The HTTP port to bind to")
}

func main() {
	cache = mapbased.NewStorage()
	tcpRun()
}

func httpRun() {
	http := httpserver.NewServer(flags.httpAddr, accounts, cache)
	if err := http.Run(); err != nil {
		log.Fatalln(err)
	}
}

func httpDevRun() {
	ttt := mapbased.TestTestStorage()
	http := httpserver.NewServer(flags.httpAddr, accounts, ttt)
	if err := http.Run(); err != nil {
		log.Fatalln(err)
	}
}

func tcpRun() {
	tcp := tcpserver.NewServer(flags.tcpAddr, cache)
	if err := tcp.Run(); err != nil {
		log.Fatalln(err)
	}
}

//TODO Заменить типы string, []string, map[string]string  на собственные алиасы этих типов
//TODO Вынести error-ы сервиса в errors.go
