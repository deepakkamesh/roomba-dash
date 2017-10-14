package dash

import (
	"fmt"
	"math"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
	ui "github.com/gizak/termui"
	"github.com/golang/glog"
)

type Dash struct {
	gui         bool
	Id          string
	roomba      *roomba.Roomba
	tmstmp      *ui.Par
	head        *ui.Par
	modeDisp    *ui.Par
	stasis      *ui.Par
	cmd         *ui.List
	irCode      *ui.Table
	movSensor   *ui.BarChart
	velSensor   *ui.BarChart
	currSensor  *ui.BarChart
	battMeter   *ui.Gauge
	dirtLvl     *ui.Gauge
	battLvl     *ui.LineChart
	battState   *ui.Table
	ocSensor    *ui.Table
	bumpSensor  *ui.Table
	wheelSensor *ui.Table
	cliffSensor *ui.Table
}

// Init initializes a a new dashboard.
func Init(gui bool, tty string, brc string) (*Dash, error) {

	m := Dash{
		Id: "dkg",
	}

	// Initialize Roomba.
	if tty != "" {
		r, err := roomba.MakeRoomba(tty, brc)
		if err != nil {
			return nil, fmt.Errorf("unable to make Roomba %v", err)
		}
		if err := r.Start(true); err != nil {
			return nil, err
		}
		m.roomba = r
	}

	// If GUI is enabled, initialize UI and Widgets.
	if gui {
		if err := ui.Init(); err != nil {
			return nil, err
		}

		m.head = ui.NewPar("")
		m.tmstmp = ui.NewPar("")
		m.modeDisp = ui.NewPar("")
		m.stasis = ui.NewPar("")
		m.cmd = ui.NewList()
		m.irCode = ui.NewTable()
		m.movSensor = ui.NewBarChart()
		m.velSensor = ui.NewBarChart()
		m.currSensor = ui.NewBarChart()
		m.battMeter = ui.NewGauge()
		m.dirtLvl = ui.NewGauge()
		m.battLvl = ui.NewLineChart()
		m.battState = ui.NewTable()
		m.ocSensor = ui.NewTable()
		m.bumpSensor = ui.NewTable()
		m.wheelSensor = ui.NewTable()
		m.cliffSensor = ui.NewTable()
		m.gui = true
	}

	return &m, nil
}

// Build creates a new dashboard layout.
func (m *Dash) Build() error {
	m.head.Text = "IRobot Create 2 - Command and Control"
	m.head.Height = 3

	m.tmstmp.Height = 3

	m.modeDisp.Height = 3

	m.stasis.Height = 3

	// Command list.
	m.cmd.Items = []string{
		"[↑] Move forward",
		"[→] Turn Right",
		"[↓] Move backward",
		"[←] Turn Right",
		"[p] Passive Mode ",
		"[f] Full Mode",
		"[s] Safe Mode",
		"[o] Stop",
		"[d] Power Down",
		"[k] Seek Dock",
	}
	m.cmd.ItemFgColor = ui.ColorYellow
	m.cmd.Height = 13
	m.cmd.BorderLabel = "Commands"
	m.cmd.Align()

	// Encoder, rotation levels.
	m.movSensor.BorderLabel = "Encoder, Rotation"
	m.movSensor.Data = []int{3231, 2223, 5234, 32, 312}
	m.movSensor.Height = 10
	m.movSensor.DataLabels = []string{"Enc(L)", "Enc(R)", "Rad(mm)", "Ang(deg)", "Dist(mm)"}
	m.movSensor.TextColor = ui.ColorGreen
	m.movSensor.BarColor = ui.ColorMagenta
	m.movSensor.NumColor = ui.ColorYellow
	m.movSensor.BarWidth = 8
	m.movSensor.Align()

	// Velocity.
	m.velSensor.BorderLabel = "Velocity (mm/s)"
	m.velSensor.Data = []int{231, 200, 100}
	m.velSensor.Height = 10
	m.velSensor.DataLabels = []string{"Total", "Right", "Left"}
	m.velSensor.TextColor = ui.ColorGreen
	m.velSensor.BarColor = ui.ColorBlue
	m.velSensor.NumColor = ui.ColorYellow
	m.velSensor.BarWidth = 5
	m.velSensor.Width = 3
	m.velSensor.Align()

	// Currents levels.
	m.currSensor.BorderLabel = "Motor Current (mAh)"
	m.currSensor.Data = []int{2223, 5234, 3223, 3122}
	m.currSensor.Height = 10
	m.currSensor.DataLabels = []string{"Left", "Right", "Main", "Side"}
	m.currSensor.TextColor = ui.ColorGreen
	m.currSensor.BarColor = ui.ColorBlue
	m.currSensor.NumColor = ui.ColorYellow
	m.currSensor.BarWidth = 5
	m.currSensor.Align()

	// Dirtlevel 0-255
	m.dirtLvl.Percent = 0
	m.dirtLvl.Height = 3
	m.dirtLvl.BorderLabel = "Dirt Level"
	m.dirtLvl.Label = "({{percent}}%) dirt detection"
	m.dirtLvl.PercentColor = ui.ColorYellow
	m.dirtLvl.BarColor = ui.ColorGreen
	m.dirtLvl.PercentColorHighlighted = ui.ColorBlack

	// Battery state gauges.
	m.battMeter.Percent = 50
	m.battMeter.Height = 3
	m.battMeter.BorderLabel = "Batt Level"
	m.battMeter.Label = "({{percent}}%) 1500/2698 mAH"
	m.battMeter.PercentColor = ui.ColorYellow
	m.battMeter.BarColor = ui.ColorGreen
	m.battMeter.PercentColorHighlighted = ui.ColorBlack

	m.battLvl.BorderLabel = "Batt Level (%)"
	m.battLvl.Data = []float64{0}
	m.battLvl.Height = 13
	//m.battLvl.Mode = "dot"
	m.battLvl.AxesColor = ui.ColorWhite
	m.battLvl.LineColor = ui.ColorGreen | ui.AttrBold

	// Battery data.
	m.battState.Rows = [][]string{
		[]string{"Batt Status", ""},
		[]string{"Temp (C)", "1023"},
		[]string{"Volts (mV)", ""},
		[]string{"Current (mA)", ""},
		[]string{"Charge", "Trickle"},
	}
	m.battState.FgColor = ui.ColorWhite
	m.battState.BgColor = ui.ColorDefault
	m.battState.TextAlign = ui.AlignCenter
	m.battState.Separator = false
	m.battState.Analysis()
	m.battState.SetSize()
	m.battState.Border = true

	// IRCode.
	m.irCode.Rows = [][]string{
		[]string{"IR Code", ""},
		[]string{"Omni", "1023"},
		[]string{"Left", ""},
		[]string{"Right", ""},
	}
	m.irCode.FgColor = ui.ColorWhite
	m.irCode.BgColor = ui.ColorDefault
	m.irCode.TextAlign = ui.AlignCenter
	m.irCode.Separator = false
	m.irCode.Border = true
	m.irCode.Analysis()
	m.irCode.SetSize()
	m.irCode.Height = 10

	// OverCurrent Data.
	m.ocSensor.Rows = [][]string{
		[]string{"Overcurrent"},
		[]string{"Right Wheel"},
		[]string{"Left Wheel"},
		[]string{"Main Brush"},
		[]string{"Side Brush"},
	}
	m.ocSensor.FgColor = ui.ColorWhite
	m.ocSensor.BgColor = ui.ColorDefault
	m.ocSensor.TextAlign = ui.AlignCenter
	m.ocSensor.Separator = false
	m.ocSensor.Analysis()
	m.ocSensor.SetSize()
	m.ocSensor.Border = true

	// Bump sensors.
	m.bumpSensor.Rows = [][]string{
		[]string{"Light Bumper", "Signal"},
		[]string{"Left", "1023"},
		[]string{"Front Left", ""},
		[]string{"Center Left", ""},
		[]string{"Center Right", ""},
		[]string{"Front Right", ""},
		[]string{"Right", ""},
	}

	m.bumpSensor.FgColor = ui.ColorWhite
	m.bumpSensor.BgColor = ui.ColorDefault
	m.bumpSensor.TextAlign = ui.AlignCenter
	m.bumpSensor.Separator = false
	m.bumpSensor.Analysis()
	m.bumpSensor.SetSize()
	m.bumpSensor.BgColors[1] = ui.ColorRed
	m.bumpSensor.BgColors[3] = ui.ColorRed
	m.bumpSensor.Border = true
	m.bumpSensor.Height = 10

	// Cliff sensors.
	m.cliffSensor.Rows = [][]string{
		[]string{"Cliff Sensor", "Signal"},
		[]string{"Left", "1023"},
		[]string{"Left Front", ""},
		[]string{"Front Right", ""},
		[]string{"Right", ""},
		[]string{"Wall", ""},
	}
	m.cliffSensor.FgColor = ui.ColorWhite
	m.cliffSensor.BgColor = ui.ColorDefault
	m.cliffSensor.TextAlign = ui.AlignCenter
	m.cliffSensor.Separator = false
	m.cliffSensor.Analysis()
	m.cliffSensor.SetSize()
	m.cliffSensor.Border = true
	m.cliffSensor.Height = 10

	// External Sensors.
	m.wheelSensor.Rows = [][]string{
		[]string{"Wheel Sensor"},
		[]string{"Right Bump"},
		[]string{"Left Bump"},
		[]string{"Right Drop"},
		[]string{"Left Drop"},
	}
	m.wheelSensor.FgColor = ui.ColorWhite
	m.wheelSensor.BgColor = ui.ColorDefault
	m.wheelSensor.TextAlign = ui.AlignCenter
	m.wheelSensor.Separator = false
	m.wheelSensor.Analysis()
	m.wheelSensor.SetSize()
	m.wheelSensor.BgColors[1] = ui.ColorRed
	m.wheelSensor.BgColors[3] = ui.ColorRed
	m.wheelSensor.Border = true
	m.wheelSensor.Height = 10

	// Align widgets.
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(9, 0, m.head),
			ui.NewCol(3, 0, m.tmstmp),
		),
		ui.NewRow(
			ui.NewCol(2, 0, m.cmd),
			ui.NewCol(3, 0, m.battMeter, m.battState, m.dirtLvl),
			ui.NewCol(2, 0, m.modeDisp, m.ocSensor, m.stasis),
			ui.NewCol(5, 0, m.battLvl),
		),
		ui.NewRow(
			ui.NewCol(3, 0, m.currSensor),
			ui.NewCol(5, 0, m.movSensor),
			ui.NewCol(4, 0, m.irCode),
		),
		ui.NewRow(
			ui.NewCol(3, 0, m.cliffSensor),
			ui.NewCol(2, 0, m.velSensor),
			ui.NewCol(3, 0, m.bumpSensor),
			ui.NewCol(2, 0, m.wheelSensor),
		),
	)
	return nil
}

const PAGEWIDTH = 120

func (m *Dash) Run() error {

	ui.Body.Width = ui.TermWidth()
	if ui.TermWidth() > PAGEWIDTH {
		ui.Body.Width = PAGEWIDTH
	}
	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("sys/kbd/<down>", func(ui.Event) {
		m.roomba.Drive(-100, 32767)
		go func() {
			t := time.Tick(500 * time.Millisecond)
			<-t
			m.roomba.Drive(0, 0)
		}()
	})

	ui.Handle("sys/kbd/<left>", func(ui.Event) {
		m.roomba.Drive(50, 1)
		go func() {
			t := time.Tick(500 * time.Millisecond)
			<-t
			m.roomba.Drive(0, 0)
		}()
	})

	ui.Handle("sys/kbd/<right>", func(ui.Event) {
		m.roomba.Drive(50, -1)
		go func() {
			t := time.Tick(500 * time.Millisecond)
			<-t
			m.roomba.Drive(0, 0)
		}()
	})

	ui.Handle("sys/kbd/<up>", func(ui.Event) {
		m.roomba.Drive(100, 32767)
		go func() {
			t := time.Tick(500 * time.Millisecond)
			<-t
			m.roomba.Drive(0, 0)
		}()
	})

	ui.Handle("/sys/kbd/p", func(ui.Event) {
		m.roomba.Passive()

	})
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		m.roomba.Safe()
	})

	ui.Handle("/sys/kbd/o", func(ui.Event) {
		m.roomba.Stop()
	})

	ui.Handle("/sys/kbd/k", func(ui.Event) {
		m.roomba.SeekDock()
	})

	ui.Handle("/sys/kbd/f", func(ui.Event) {
		m.roomba.Full()
	})

	ui.Handle("/sys/kbd/d", func(ui.Event) {
		if m.roomba != nil {
			m.roomba.Power()
		}
		glog.Flush()
		ui.StopLoop()
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		ui.Body.Width = ui.TermWidth()
		if ui.TermWidth() > PAGEWIDTH {
			ui.Body.Width = PAGEWIDTH
		}
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)

	})

	ui.Loop()
	ui.Close()
	return nil
}

func (m *Dash) Update() error {

	if m.roomba == nil {
		return fmt.Errorf("roomba not initialized")
	}

	t := time.NewTicker(1000 * time.Millisecond)
	sg := []byte{constants.SENSOR_GROUP_6, constants.SENSOR_GROUP_101}
	pg := [][]byte{constants.PACKET_GROUP_6, constants.PACKET_GROUP_101}

	for {
		<-t.C
		fmt.Printf("%s\n", time.Now().Format("2006/01/02 150405"))

		// Iterate through the packet groups. Sensor group 100 does not work as advertised.
		// Use sensor group, 6 and 101 instead.
		for grp := 0; grp < 2; grp++ {
			d, e := m.roomba.Sensors(sg[grp])
			if e != nil {
				return e
			}

			i := byte(0)
			for _, p := range pg[grp] {
				pktL := constants.SENSOR_PACKET_LENGTH[p]

				if pktL == 1 {
					fmt.Printf("%25s:  %v \n", constants.SENSORS_NAME[p], d[i])
				}
				if pktL == 2 {

					fmt.Printf("%25s:  %v \n", constants.SENSORS_NAME[p], int16(d[i])<<8|int16(d[i+1]))
				}
				i = i + pktL
			}
		}
	}
}

func (m *Dash) UpdateGUI() error {
	if m.roomba == nil {
		return fmt.Errorf("roomba not initialized")
	}

	battMaxPlot := 0
	var battCap uint16

	t := time.NewTicker(300 * time.Millisecond)
	sg := []byte{constants.SENSOR_GROUP_6, constants.SENSOR_GROUP_101}
	pg := [][]byte{constants.PACKET_GROUP_6, constants.PACKET_GROUP_101}

	for {

		<-t.C
		m.tmstmp.Text = fmt.Sprintf(time.Now().Format("2006/01/02 - 15-04-05"))

		// Iterate through the packet groups. Sensor group 100 does not work as advertised.
		// Use sensor group, 6 and 101 instead.
		for grp := 0; grp < 2; grp++ {
			// TODO: Need a timeout to return from function.
			d, e := m.roomba.Sensors(sg[grp])
			if e != nil {
				// TODO: Log Error and continue.
				glog.Errorf("Error reading sensors %v", e)
				continue
			}
			i := byte(0)

			for _, p := range pg[grp] {
				switch p {
				case constants.SENSOR_TEMPERATURE:
					m.battState.Rows[1][1] = fmt.Sprintf("%d", d[i])

				case constants.SENSOR_VOLTAGE:
					m.battState.Rows[2][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_CURRENT:
					m.battState.Rows[3][1] = fmt.Sprintf("%d", int16(d[i])<<8|int16(d[i+1]))

				case constants.SENSOR_BATTERY_CAPACITY:
					battCap = uint16(d[i])<<8 | uint16(d[i+1])

				case constants.SENSOR_CHARGING:
					ch, ok := constants.CHARGING_STATE[d[i]]
					if !ok {
						m.battState.Rows[4][1] = "Unknown"
					}
					m.battState.Rows[4][1] = ch

				case constants.SENSOR_BATTERY_CHARGE:
					bc := uint16(d[i])<<8 | uint16(d[i+1])
					battMaxPlot = battMaxPlot + 1
					perc := float64(0)
					if battCap > 0 {
						perc = float64(bc) * 100 / float64(battCap)
					}
					if battMaxPlot%20 == 0 {
						m.battLvl.Data = append(m.battLvl.Data, toFixed(perc, 2))
						if battMaxPlot%1200 == 0 {
							m.battLvl.Data = []float64{}
						}
					}
					m.battMeter.Percent = int(perc)
					m.battMeter.Label = fmt.Sprintf(" %.2f%% %d/%d mAh", perc, bc, battCap)

				case constants.SENSOR_IR_OMNI:
					n, ok := constants.IR_CODE_NAMES[d[i]]
					if !ok {
						n = "unknown"
					}
					m.irCode.Rows[1][1] = n

				case constants.SENSOR_IR_LEFT:
					n, ok := constants.IR_CODE_NAMES[d[i]]
					if !ok {
						n = "unknown"
					}
					m.irCode.Rows[2][1] = n

				case constants.SENSOR_IR_RIGHT:
					n, ok := constants.IR_CODE_NAMES[d[i]]
					if !ok {
						n = "unknown"
					}
					m.irCode.Rows[3][1] = n

				case constants.SENSOR_OI_MODE:
					mode, ok := constants.OI_MODE[d[i]]
					if !ok {
						m.modeDisp.Text = "Unknown Mode"
						break
					}
					m.modeDisp.Text = "Mode:" + mode

				case constants.SENSOR_CLIFF_LEFT_SIGNAL:
					m.cliffSensor.Rows[1][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_CLIFF_FRONT_LEFT_SIGNAL:
					m.cliffSensor.Rows[2][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_CLIFF_FRONT_RIGHT_SIGNAL:
					m.cliffSensor.Rows[3][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_CLIFF_RIGHT_SIGNAL:
					m.cliffSensor.Rows[4][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_WALL_SIGNAL:
					m.cliffSensor.Rows[5][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_CLIFF_LEFT:
					if d[i] == 1 {
						m.cliffSensor.BgColors[1] = ui.ColorRed
						break
					}
					m.cliffSensor.BgColors[1] = ui.ColorDefault

				case constants.SENSOR_CLIFF_FRONT_LEFT:
					if d[i] == 1 {
						m.cliffSensor.BgColors[2] = ui.ColorRed
						break
					}
					m.cliffSensor.BgColors[2] = ui.ColorDefault

				case constants.SENSOR_CLIFF_FRONT_RIGHT:
					if d[i] == 1 {
						m.cliffSensor.BgColors[3] = ui.ColorRed
						break
					}
					m.cliffSensor.BgColors[3] = ui.ColorDefault

				case constants.SENSOR_CLIFF_RIGHT:
					if d[i] == 1 {
						m.cliffSensor.BgColors[4] = ui.ColorRed
						break
					}
					m.cliffSensor.BgColors[4] = ui.ColorDefault

				case constants.SENSOR_WALL:
					if d[i] == 1 {
						m.cliffSensor.BgColors[5] = ui.ColorRed
						break
					}
					m.cliffSensor.BgColors[5] = ui.ColorDefault

				case constants.SENSOR_WHEEL_OVERCURRENT:
					m.ocSensor.BgColors[1] = ui.ColorDefault
					m.ocSensor.BgColors[2] = ui.ColorDefault
					m.ocSensor.BgColors[3] = ui.ColorDefault
					m.ocSensor.BgColors[4] = ui.ColorDefault
					switch {
					case d[i]&1 > 0:
						m.ocSensor.BgColors[4] = ui.ColorRed
					case d[i]&4 > 0:
						m.ocSensor.BgColors[3] = ui.ColorRed
					case d[i]&8 > 0:
						m.ocSensor.BgColors[2] = ui.ColorRed
					case d[i]&16 > 0:
						m.ocSensor.BgColors[1] = ui.ColorRed
					}

				case constants.SENSOR_BUMP_WHEELS_DROPS:
					idx := 1
					for offset := byte(1); offset <= 8; offset = offset << 1 {
						m.wheelSensor.BgColors[idx] = ui.ColorDefault
						if d[i]&offset > 0 {
							m.wheelSensor.BgColors[idx] = ui.ColorRed
						}
						idx++
					}

				case constants.SENSOR_BUMPER:
					idx := 1
					for offset := byte(1); offset <= 32; offset = offset << 1 {
						m.bumpSensor.BgColors[idx] = ui.ColorDefault
						if d[i]&offset > 0 {
							m.bumpSensor.BgColors[idx] = ui.ColorRed
						}
						idx++
					}

				case constants.SENSOR_BUMP_LEFT:
					m.bumpSensor.Rows[1][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_BUMP_FRONT_LEFT:
					m.bumpSensor.Rows[2][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_BUMP_CENTER_LEFT:
					m.bumpSensor.Rows[3][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_BUMP_CENTER_RIGHT:
					m.bumpSensor.Rows[4][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_BUMP_FRONT_RIGHT:
					m.bumpSensor.Rows[5][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_BUMP_RIGHT:
					m.bumpSensor.Rows[6][1] = fmt.Sprintf("%d", uint16(d[i])<<8|uint16(d[i+1]))

				case constants.SENSOR_LEFT_MOTOR_CURRENT:
					m.currSensor.Data[0] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_RIGHT_MOTOR_CURRENT:
					m.currSensor.Data[1] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_MAIN_BRUSH_MOTOR_CURRENT:
					m.currSensor.Data[2] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_SIDE_BRUSH_MOTOR_CURRENT:
					m.currSensor.Data[3] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_REQUESTED_VELOCITY:
					m.velSensor.Data[0] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_REQUESTED_RIGHT_VELOCITY:
					m.velSensor.Data[1] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_REQUESTED_LEFT_VELOCITY:
					m.velSensor.Data[2] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_STASIS:
					switch d[i] {
					case 0:
						m.stasis.Text = "Stasis:BWD/TURN"
					case 1:
						m.stasis.Text = "Stasis:FWD"
					}

				case constants.SENSOR_DIRT_DETECT:
					m.dirtLvl.Percent = int(d[i] * 100 / 255)

				case constants.SENSOR_ANGLE:
					m.movSensor.Data[3] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_DISTANCE:
					m.movSensor.Data[4] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_REQUESTED_RADIUS:
					m.movSensor.Data[2] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_LEFT_ENCODER:
					m.movSensor.Data[0] = int(int16(d[i])<<8 | int16(d[i+1]))

				case constants.SENSOR_RIGHT_ENCODER:
					m.movSensor.Data[1] = int(int16(d[i])<<8 | int16(d[i+1]))
				}
				i = i + constants.SENSOR_PACKET_LENGTH[p]
			}
			ui.Body.Align()
			ui.Render(m.velSensor)
			ui.Render(ui.Body)
		}
	}
	return nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
