package main

import (
	"fmt"

	"github.com/hood-chat/core"
	"github.com/hood-chat/core/entity"
	"github.com/hood-chat/core/repo"
	"github.com/hood-chat/core/store"

	logging "github.com/ipfs/go-log"
)

func init() {
	fmt.Println("Hello! init() function")
}

// Main function
func main() {
	err := logging.SetLogLevel("*", "DEBUG")
	if err != nil {
		panic(err)
	}
	s, err := store.NewStore("./data")
	if err != nil {
		panic(err)
	}
	rIdentity := repo.NewIdentityRepo(s)
	id, err := rIdentity.Get()
	if err != nil {
		id, err = entity.CreateIdentity("bootstraper")
		if err != nil {
			panic(err)
		}
	}
	opt := core.DefaultOption()
	opt.SetIdentity(&id)
	hb := core.DefaultRoutedHost{}
	if err != nil {
		panic(err)
	}

	_, err = hb.Create(opt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Welcome to main() function")

	select {} // block forever
}
