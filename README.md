# roomba-dash

## Install/Usage Instructions.
This is meant to be a go package for use with other binaries. A simple implementation is provided in bin/main.go. It assumes the roomba is connected to a raspberry pi device and also provides for flipping BRC to keep the roomba from sleeping. An even simpler implementation is provided below without the BRC code. It should work from any machine that is connected to roomba via the provided cable. The cable is a usb to serial cable and it sets up a serial device on your machine. You can find out which serial port by tail -f /var/log/message as you plug in the cable. Pass that serial port (usually /dev/ttyS0 or something similar on linux) as a flag below.

```
func main() {
	ui := flag.Bool("ui", false, "Disable UI")
	tty := flag.String("tty", "/dev/ttyS0", "Serial TTY device") 
	flag.Parse()

        // Initialize the dashboard.
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
```
