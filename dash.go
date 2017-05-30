package dash

import (
	ui "github.com/gizak/termui"
)

type Dash struct {
	Id string
}

// Init initializes a a new dashboard.
func Init() (*Dash, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	return &Dash{
		Id: "dd",
	}, nil
}

// Build creates a new dashboard layout.

func (m *Dash) Build() error {
	cmds := ui.NewPar("IRobot Create 2 Command and Control")
	cmds.Height = 3

	// Command list.
	cmd := ui.NewList()
	cmd.Items = []string{
		"[↑] Move forward",
		"[→] Turn Right",
		"[↓] Move backward",
		"[←] Turn Right",
		"[p] Passive Mode ",
		"[f] Full Mode",
		"[s] Safe Mode",
		"[c] ",
	}
	cmd.ItemFgColor = ui.ColorYellow
	cmd.Height = 10
	cmd.BorderLabel = "Commands"

	// Encoder, rotation levels.
	movSensor := ui.NewBarChart()
	movData := []int{3231, 2223, 5234, 32, 312}
	movlabels := []string{"Enc(L)", "Enc(R)", "Rad(mm)", "Ang(deg)", "Dist(mm)"}
	movSensor.BorderLabel = "Encoder, Rotation"
	movSensor.Data = movData
	movSensor.Height = 10
	movSensor.DataLabels = movlabels
	movSensor.TextColor = ui.ColorGreen
	movSensor.BarColor = ui.ColorMagenta
	movSensor.NumColor = ui.ColorYellow
	movSensor.BarWidth = 8
	movSensor.Align()

	// Velocity.
	velSensor := ui.NewBarChart()
	velData := []int{231, 200, 100}
	vellabels := []string{"Total", "Right", "Left"}
	velSensor.BorderLabel = "Velocity (mm/s)"
	velSensor.Data = velData
	velSensor.Height = 10
	velSensor.DataLabels = vellabels
	velSensor.TextColor = ui.ColorGreen
	velSensor.BarColor = ui.ColorBlue
	velSensor.NumColor = ui.ColorYellow
	velSensor.BarWidth = 5
	velSensor.Align()

	// Currents levels.
	currSensor := ui.NewBarChart()
	data := []int{3231, 2223, 5234, 3223, 3122}
	bclabels := []string{"Batt", "Left", "Right", "Main", "Side"}
	currSensor.BorderLabel = "Current Levels (mAh)"
	currSensor.Data = data
	currSensor.Height = 10
	currSensor.DataLabels = bclabels
	currSensor.TextColor = ui.ColorGreen
	currSensor.BarColor = ui.ColorBlue
	currSensor.NumColor = ui.ColorYellow
	currSensor.BarWidth = 5
	currSensor.Align()

	// Battery state gauges.
	battPer := ui.NewGauge()
	battPer.Percent = 50
	battPer.Height = 3
	battPer.BorderLabel = "Batt Level"
	battPer.Label = "({{percent}}%) 1500/2698 mAH"
	battPer.PercentColor = ui.ColorYellow
	battPer.BarColor = ui.ColorGreen
	battPer.PercentColorHighlighted = ui.ColorBlack

	batt := ui.NewLineChart()
	batt.BorderLabel = "Batt Level (mAh)"
	batt.Data = []float64{3.2, 3.3, 3.6, 2.3, 4, 5, 9, 2}
	batt.Height = 10
	batt.Mode = "dot"
	batt.AxesColor = ui.ColorWhite
	batt.LineColor = ui.ColorGreen | ui.AttrBold

	// Battery data.
	battStateData := [][]string{
		[]string{"Batt Status", ""},
		[]string{"Temp (C)", "1023"},
		[]string{"Volts (mV)", ""},
		[]string{"Current (mA)", ""},
		[]string{"Charge", "Trickle"},
	}
	battState := ui.NewTable()
	battState.Rows = battStateData
	battState.FgColor = ui.ColorWhite
	battState.BgColor = ui.ColorDefault
	battState.TextAlign = ui.AlignCenter
	battState.Separator = false
	battState.Analysis()
	battState.SetSize()
	battState.Border = true

	// OverCurrent Data.
	ocData := [][]string{
		[]string{"Overcurrent"},
		[]string{"Right Wheel"},
		[]string{"Left Wheel"},
		[]string{"Main Brush"},
		[]string{"Side Brush"},
	}
	ocSensor := ui.NewTable()
	ocSensor.Rows = ocData
	ocSensor.FgColor = ui.ColorWhite
	ocSensor.BgColor = ui.ColorDefault
	ocSensor.TextAlign = ui.AlignCenter
	ocSensor.Separator = false
	ocSensor.Analysis()
	ocSensor.SetSize()
	ocSensor.Border = true

	// Bump sensors.
	bumpData := [][]string{
		[]string{"Light Bumper", "Signal"},
		[]string{"Left", "1023"},
		[]string{"Front Left", ""},
		[]string{"Center Left", ""},
		[]string{"Center Right", ""},
		[]string{"Front Right", ""},
		[]string{"Right", ""},
	}
	bumpSensor := ui.NewTable()
	bumpSensor.Rows = bumpData
	bumpSensor.FgColor = ui.ColorWhite
	bumpSensor.BgColor = ui.ColorDefault
	bumpSensor.TextAlign = ui.AlignCenter
	bumpSensor.Separator = false
	bumpSensor.Analysis()
	bumpSensor.SetSize()
	bumpSensor.BgColors[1] = ui.ColorRed
	bumpSensor.BgColors[3] = ui.ColorRed
	bumpSensor.Border = true

	// Cliff sensors.
	cliffData := [][]string{
		[]string{"Cliff Sensor", "Signal"},
		[]string{"Left", "1023"},
		[]string{"Left Front", ""},
		[]string{"Right", ""},
		[]string{"Front Right", ""},
	}
	cliffSensor := ui.NewTable()
	cliffSensor.Rows = cliffData
	cliffSensor.FgColor = ui.ColorWhite
	cliffSensor.BgColor = ui.ColorDefault
	cliffSensor.TextAlign = ui.AlignCenter
	cliffSensor.Separator = false
	cliffSensor.Analysis()
	cliffSensor.SetSize()
	cliffSensor.BgColors[1] = ui.ColorRed
	cliffSensor.BgColors[3] = ui.ColorRed
	cliffSensor.Border = true

	// External Sensors.
	rows2 := [][]string{
		[]string{"Wheel Sensor", ""},
		[]string{"Right Drop", ""},
		[]string{"Right Drop", ""},
		[]string{"Left Bump", ""},
		[]string{"Right Bump", ""},
	}
	wheelSensor := ui.NewTable()
	wheelSensor.Rows = rows2
	wheelSensor.FgColor = ui.ColorWhite
	wheelSensor.BgColor = ui.ColorDefault
	wheelSensor.TextAlign = ui.AlignCenter
	wheelSensor.Separator = false
	wheelSensor.Analysis()
	wheelSensor.SetSize()
	wheelSensor.BgColors[1] = ui.ColorRed
	wheelSensor.BgColors[3] = ui.ColorRed
	wheelSensor.Border = true

	// Align widgets.
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, cmds),
		),
		ui.NewRow(
			ui.NewCol(3, 0, cmd),
			ui.NewCol(3, 0, battPer, battState),
			ui.NewCol(4, 0, batt),
			ui.NewCol(2, 0, ocSensor),
		),
		ui.NewRow(
			ui.NewCol(4, 0, currSensor),
			ui.NewCol(5, 0, movSensor),
		),
		ui.NewRow(
			ui.NewCol(3, 0, cliffSensor),
			ui.NewCol(3, 0, velSensor),
			ui.NewCol(3, 0, bumpSensor),
			ui.NewCol(3, 0, wheelSensor),
		),
	)
	return nil
}

func (m *Dash) Run() error {
	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		ui.Body.Align()
		ui.Render(ui.Body)
		cmds.Text = "sss"
	})

	ui.Loop()
	ui.Close()
	return nil

}
