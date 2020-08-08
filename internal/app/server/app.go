package server

import (
	"fmt"
)

func Run(arg Args) error {
	fmt.Printf("%+v", arg)
	return nil
}
