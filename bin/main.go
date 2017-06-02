package main

import (
	"flag"
	"fmt"

	"github.com/deepakkamesh/roomba-dash"
)

func main() {
	ui := flag.Bool("ui", false, "Disable UI")
	tty := flag.String("tty", "", "Serial TTY device")
	brc := flag.String("brc", "LCD-D23", "BRC GPIO port")
	flag.Parse()

	d, err := dash.Init(*ui, *tty, *brc)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	if !(*ui) {
		d.Update()
		return
	}

	d.Build()
	go d.UpdateGUI()
	d.Run()

}
