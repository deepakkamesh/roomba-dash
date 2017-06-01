package dash

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
	"github.com/xa4a/go-roomba"
	"github.com/xa4a/go-roomba/constants"
)

type Dash struct {
	gui         bool
	Id          string
	roomba      *roomba.Roomba
	tmstmp      *ui.Par
	head        *ui.Par
	modeDisp    *ui.Par
	cmd         *ui.List
	irCode      *ui.Table
	movSensor   *ui.BarChart
	velSensor   *ui.BarChart
	currSensor  *ui.BarChart
	battMeter   *ui.Gauge
	battLvl     *ui.LineChart
	battState   *ui.Table
	ocSensor    *ui.Table
	bumpSensor  *ui.Table
	wheelSensor *ui.Table
	cliffSensor *ui.Table
}

// Init initializes a a new dashboard.
func Init(gui bool, tty string) (*Dash, error) {

	m := Dash{
		Id: "dkg",
	}

	// Initialize Roomba.
	if tty != "" {
		r, err := roomba.MakeRoomba(tty)
		if err != nil {
			return nil, fmt.Errorf("unable to make Roomba %v", err)
		}
		r.Start()
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
		m.cmd = ui.NewList()
		m.irCode = ui.NewTable()
		m.movSensor = ui.NewBarChart()
		m.velSensor = ui.NewBarChart()
		m.currSensor = ui.NewBarChart()
		m.battMeter = ui.NewGauge()
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
	}
	m.cmd.ItemFgColor = ui.ColorYellow
	m.cmd.Height = 12
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

	// Battery state gauges.
	m.battMeter.Percent = 50
	m.battMeter.Height = 3
	m.battMeter.BorderLabel = "Batt Level"
	m.battMeter.Label = "({{percent}}%) 1500/2698 mAH"
	m.battMeter.PercentColor = ui.ColorYellow
	m.battMeter.BarColor = ui.ColorGreen
	m.battMeter.PercentColorHighlighted = ui.ColorBlack

	m.battLvl.BorderLabel = "Batt Level (mAh)"
	m.battLvl.Data = []float64{}
	m.battLvl.Height = 12
	m.battLvl.Mode = "dot"
	m.battLvl.AxesColor = ui.ColorWhite
	m.battLvl.LineColor = ui.ColorGreen | ui.AttrBold

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

	// External Sensors.
	m.wheelSensor.Rows = [][]string{
		[]string{"Wheel Sensor"},
		[]string{"Right Drop"},
		[]string{"Right Drop"},
		[]string{"Left Bump"},
		[]string{"Right Bump"},
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

	// Align widgets.
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(9, 0, m.head),
			ui.NewCol(3, 0, m.tmstmp),
		),
		ui.NewRow(
			ui.NewCol(2, 0, m.cmd),
			ui.NewCol(3, 0, m.battMeter, m.battState),
			ui.NewCol(2, 0, m.modeDisp, m.ocSensor),
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

	ui.Handle("/sys/kbd/p", func(ui.Event) {
		m.roomba.Passive()

	})
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		m.roomba.Safe()
	})

	ui.Handle("/sys/kbd/o", func(ui.Event) {
		m.roomba.Stop()
	})

	ui.Handle("/sys/kbd/f", func(ui.Event) {
		m.roomba.Full()
	})

	ui.Handle("/sys/kbd/d", func(ui.Event) {
		m.roomba.Power()
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
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
			d, e := m.roomba.Sensors(sg[grp])
			if e != nil {
				return e
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
					battMaxPlot = battMaxPlot + 1
					bc := uint16(d[i])<<8 | uint16(d[i+1])
					m.battLvl.Data = append(m.battLvl.Data, float64(bc))
					if battMaxPlot > 50 {
						m.battLvl.Data = []float64{}
						battMaxPlot = 0
					}
					perc := float64(0)
					if battCap > 0 {
						perc = float64(bc) * 100 / float64(battCap)
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
