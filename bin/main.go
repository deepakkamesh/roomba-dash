package main

import (
	"fmt"

	"github.com/deepakkamesh/roomba-dash"
)

func main() {
	d, _ := dash.Init()
	d.Build()
	d.Run()
	fmt.Println("dont")
}
