package main

import (
	"fmt"
	"net/http"

	"github.com/fulldump/box"
	"github.com/fulldump/goconfig"

	"instantlogs/api"
	"instantlogs/service"
)

type Config struct {
	Addr       string
	StaticsDir string
}

func main() {

	fmt.Println(banner)

	c := &Config{
		Addr: ":8080",
	}
	goconfig.Read(&c)

	s := service.NewService()

	a := api.NewApi(s, c.StaticsDir)

	server := &http.Server{
		Addr:    c.Addr,
		Handler: box.Box2Http(a),
	}
	fmt.Printf("Listening on %s\n", server.Addr)
	server.ListenAndServe()
}

const banner = `
________             _____              ___________                       
____  _/_______________  /______ _________  /___  / _____________ ________
 __  / __  __ \_  ___/  __/  __ '/_  __ \  __/_  /  _  __ \_  __ '/_  ___/
__/ /  _  / / /(__  )/ /_ / /_/ /_  / / / /_ _  /___/ /_/ /  /_/ /_(__  ) 
/___/  /_/ /_//____/ \__/ \__,_/ /_/ /_/\__/ /_____/\____/_\__, / /____/  
                                                          /____/          
`
