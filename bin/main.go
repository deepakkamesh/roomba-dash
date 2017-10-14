package main

import (
	"flag"
	"fmt"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"

	"github.com/deepakkamesh/roomba-dash"
)

func main() {
	ui := flag.Bool("ui", false, "Disable UI")
	tty := flag.String("tty", "", "Serial TTY device")
	brc := flag.String("brc", "7", "BRC GPIO port")
	flag.Parse()

	// Need GPIO control for BRC.
	pi := raspi.NewAdaptor()
	if err := pi.Connect(); err != nil {
		panic(err)
	}
	p := gpio.NewDirectPinDriver(pi, *brc)
	if err := p.Start(); err != nil {
		panic(err)
	}

	d, err := dash.Init(*ui, *tty, nil)
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
