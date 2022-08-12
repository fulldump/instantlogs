package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fulldump/box"
	"github.com/fulldump/goconfig"

	"instantlogs/api"
	"instantlogs/blocks"
	"instantlogs/blocks/bigblock"
	"instantlogs/blocks/blockchain"
	"instantlogs/service"
)

type Config struct {
	Addr       string
	StaticsDir string
}

func main() {

	fmt.Print(banner)

	c := &Config{
		Addr: ":8080",
	}
	goconfig.Read(&c)

	bc := blockchain.New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, 10*1024*1024))
	})

	blockCounter := 0
	bc.OnBlockCompleted(func(block blocks.Blocker) {
		blockCounter++
		fmt.Println("New block", blockCounter)
	})

	s := service.NewService(bc)

	a := api.NewApi(s, c.StaticsDir)

	server := &http.Server{
		Addr:    c.Addr,
		Handler: box.Box2Http(a),
	}
	fmt.Fprintf(os.Stderr, "Listening on %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

const banner = `
________             _____              ___________                       
____  _/_______________  /______ _________  /___  / _____________ ________
 __  / __  __ \_  ___/  __/  __ '/_  __ \  __/_  /  _  __ \_  __ '/_  ___/
__/ /  _  / / /(__  )/ /_ / /_/ /_  / / / /_ _  /___/ /_/ /  /_/ /_(__  ) 
/___/  /_/ /_//____/ \__/ \__,_/ /_/ /_/\__/ /_____/\____/_\__, / /____/  
                                                          /____/          
`
