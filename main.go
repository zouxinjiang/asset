package main

import (
	"github.com/zouxinjiang/axes/cmd"
	"github.com/zouxinjiang/axes/pkg/log"
)

func main() {
	log.SetShowLevel(log.Lvl_All)
	err := cmd.Execute()
	if err != nil {
		log.Error(err)
	}
}
