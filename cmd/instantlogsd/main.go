package main

import (
	"fmt"
	"github.com/fulldump/box"
	"github.com/fulldump/goconfig"
	"instantlogs/api"
	"instantlogs/service"
	"net/http"
)

type Config struct {
	Addr string
}

func main() {

	c := &Config{
		Addr: ":8080",
	}
	goconfig.Read(&c)

	s := service.NewService()

	a := api.NewApi(s)

	server := &http.Server{
		Addr:    c.Addr,
		Handler: box.Box2Http(a),
	}
	fmt.Printf("Listening on %s\n", server.Addr)
	server.ListenAndServe()
}
