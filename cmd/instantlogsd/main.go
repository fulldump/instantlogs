package main

import (
	"fmt"
	"log"
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
	Addr       string `usage:"HTTP address to expose the service"`
	StaticsDir string `usage:"Statics dir, if empty embedded statics will be used"`
	BlockSize  int    `usage:"BlockSize (in MiB)"`
	BlockNum   int    `usage:"Number of blocks"`
}

func main() {

	fmt.Print(banner)

	c := &Config{
		Addr:      ":8080",
		BlockSize: 20,
		BlockNum:  5,
	}
	goconfig.Read(&c)

	bc := blockchain.New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, c.BlockSize*1024*1024))
	})
	bc.MaxBlocks = c.BlockNum

	// Setup some logs
	blockCounter := 0
	bc.OnBlockCompleted(func(block blocks.Blocker) {
		blockCounter++
		log.Println("New block", blockCounter)
	})
	bc.OnBlockDiscarded(func(block blocks.Blocker) {
		log.Println("Block discarded")
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
