package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func main() {
	Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	cwd, err := os.Getwd()
	if err != nil {
		Logger.Fatalf("Failed to get current working directory: %v", err)
	}

	path, err := dialog.
		File().
		Title("Open BIN file").
		Filter("Binary files", "bin").
		SetStartDir(cwd).
		Load()
	if err != nil {
		Logger.Fatalf("Failed to select file: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		Logger.Fatalf("Failed to read file: %v", err)
	}

	Logger.Printf("Loaded %d bytes from %s\n", len(data), path)

	MainApp := app.New()
	MainApp.Settings().SetTheme(theme.DarkTheme())
	Wnd := MainApp.NewWindow("NANDPU Simulator")

	Wnd.Resize(fyne.NewSize(800, 600))
	Wnd.SetFixedSize(true)

	nandpu := NewNANDPU(data)
	var updateGUIValues func()

	running := false
	stepNum := 0
	speed := binding.NewFloat()
	speed.Set(50)

	createFixedLabel := func() (*widget.Label, *fyne.Container) {
		value := widget.NewLabel("0x0000")
		return value, container.NewWithoutLayout(value)
	}

	pcLabel, pcLabelContainer := createFixedLabel()
	spLabel, spLabelContainer := createFixedLabel()
	incLabel, incLabelContainer := createFixedLabel()

	aLabel, aLabelContainer := createFixedLabel()
	bLabel, bLabelContainer := createFixedLabel()
	cLabel, cLabelContainer := createFixedLabel()
	dLabel, dLabelContainer := createFixedLabel()

	mLabel, mLabelContainer := createFixedLabel()
	xyLabel, xyLabelContainer := createFixedLabel()
	jLabel, jLabelContainer := createFixedLabel()

	m1Label, m1LabelContainer := createFixedLabel()
	m2Label, m2LabelContainer := createFixedLabel()

	xLabel, xLabelContainer := createFixedLabel()
	yLabel, yLabelContainer := createFixedLabel()

	j1Label, j1LabelContainer := createFixedLabel()
	j2Label, j2LabelContainer := createFixedLabel()

	zeroLabel, zeroLabelContainer := createFixedLabel()
	carryLabel, carryLabelContainer := createFixedLabel()
	signLabel, signLabelContainer := createFixedLabel()
	lessThanLabel, lessThanLabelContainer := createFixedLabel()

	stepNumLabel := widget.NewLabel("0")

	var runBtn *widget.Button
	var stepBtn *widget.Button
	var resetBtn *widget.Button

	runBtn = widget.NewButton("Run", func() {
		if running {
			fmt.Println("Stop button clicked")
			running = false
		} else {
			fmt.Println("Run button clicked")
			running = true
		}
		updateGUIValues()
	})
	stepBtn = widget.NewButton("Step", func() {
		fmt.Println("Step button clicked")
		nandpu.Step()
		stepNum += 1
		updateGUIValues()
	})
	resetBtn = widget.NewButton("Reset", func() {
		fmt.Println("Reset button clicked")
		nandpu = NewNANDPU(data)
		stepNum = 0
		updateGUIValues()
	})

	stringSpeed := binding.NewString()
	speedVal, _ := speed.Get()
	stringSpeed.Set(fmt.Sprintf("%.1fms", speedVal))
	speedLabel := widget.NewLabelWithData(stringSpeed)
	speedSlider := widget.NewSliderWithData(0.1, 150.0, speed)
	speedSlider.Step = 0.1

	speedSliderContainer := container.NewGridWrap(fyne.NewSize(200, 40), speedSlider)

	btnRow := container.NewHBox(
		runBtn, stepBtn, resetBtn, stepNumLabel, speedSliderContainer, speedLabel,
	)

	speed.AddListener(binding.NewDataListener(func() {
		speedVal, _ := speed.Get()
		stringSpeed.Set(fmt.Sprintf("%.1fms", speedVal))
	}))

	regRow1 := container.NewHBox(
		widget.NewSeparator(),
		widget.NewLabel("PC"), pcLabelContainer, widget.NewSeparator(),
		widget.NewLabel("SP"), spLabelContainer, widget.NewSeparator(),
		widget.NewLabel("INC"), incLabelContainer, widget.NewSeparator(),
	)

	regRow2 := container.NewHBox(
		widget.NewSeparator(),
		widget.NewLabel("A"), aLabelContainer, widget.NewSeparator(),
		widget.NewLabel("B"), bLabelContainer, widget.NewSeparator(),
		widget.NewLabel("C"), cLabelContainer, widget.NewSeparator(),
		widget.NewLabel("D"), dLabelContainer, widget.NewSeparator(),
	)

	regRow3M := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("M"), mLabelContainer,
			layout.NewSpacer(),
		),
		container.NewHBox(
			widget.NewLabel("M1"), m1LabelContainer, widget.NewSeparator(),
			widget.NewLabel("M2"), m2LabelContainer,
		),
	)

	regRow3XY := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("XY"), xyLabelContainer,
			layout.NewSpacer(),
		),
		container.NewHBox(
			widget.NewLabel("X"), xLabelContainer, widget.NewSeparator(),
			widget.NewLabel("Y"), yLabelContainer,
		),
	)

	regRow3J := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("J"), jLabelContainer,
			layout.NewSpacer(),
		),
		container.NewHBox(
			widget.NewLabel("J1"), j1LabelContainer, widget.NewSeparator(),
			widget.NewLabel("J2"), j2LabelContainer,
		),
	)

	regRow3 := container.NewHBox(
		widget.NewSeparator(),
		regRow3M, widget.NewSeparator(),
		regRow3XY, widget.NewSeparator(),
		regRow3J, widget.NewSeparator(),
	)

	regRow4 := container.NewHBox(
		widget.NewSeparator(),
		widget.NewLabel("Z"), zeroLabelContainer, widget.NewSeparator(),
		widget.NewLabel("C"), carryLabelContainer, widget.NewSeparator(),
		widget.NewLabel("S"), signLabelContainer, widget.NewSeparator(),
		widget.NewLabel("LT"), lessThanLabelContainer, widget.NewSeparator(),
	)

	regContainer := container.NewVBox(
		btnRow,
		widget.NewSeparator(),
		regRow1,
		widget.NewSeparator(),
		regRow2,
		widget.NewSeparator(),
		regRow3,
		widget.NewSeparator(),
		regRow4,
		widget.NewSeparator(),
	)

	createMemoryList := func(mem *MemMap) *widget.List {
		const rowSize = 16
		const totalCells = 65536 // full 16-bit address space

		// total rows is totalCells / rowSize
		rowCount := totalCells / rowSize

		return widget.NewList(
			func() int {
				return rowCount
			},
			func() fyne.CanvasObject {
				// Create a row with 16 labels
				row := make([]fyne.CanvasObject, rowSize)
				for i := range row {
					row[i] = widget.NewLabel("00")
				}
				return container.NewHBox(row...)
			},
			func(id widget.ListItemID, item fyne.CanvasObject) {
				row := item.(*fyne.Container)
				for i := 0; i < rowSize; i++ {
					addr := uint16(id*rowSize + i)
					byteValue := mem.Read(addr)
					label := row.Objects[i].(*widget.Label)
					label.SetText(fmt.Sprintf("%02X", byteValue))
				}
			},
		)
	}

	memList := createMemoryList(&nandpu.Mem)

	mainContainer := container.NewBorder(
		regContainer, nil, nil, nil,
		memList,
	)

	content := container.NewStack(mainContainer)
	Wnd.SetContent(content)

	updateGUIValues = func() {
		pcLabel.SetText(fmt.Sprintf("0x%04X", nandpu.PC.val))
		spLabel.SetText(fmt.Sprintf("0x%04X", nandpu.SP.val))
		incLabel.SetText(fmt.Sprintf("0x%04X", nandpu.INC.val))

		aLabel.SetText(fmt.Sprintf("0x%02X", nandpu.RegA.val))
		bLabel.SetText(fmt.Sprintf("0x%02X", nandpu.RegB.val))
		cLabel.SetText(fmt.Sprintf("0x%02X", nandpu.RegC.val))
		dLabel.SetText(fmt.Sprintf("0x%02X", nandpu.RegD.val))

		mLabel.SetText(fmt.Sprintf("0x%04X", nandpu.RegM.val))
		xyLabel.SetText(fmt.Sprintf("0x%04X", nandpu.RegXY.val))
		jLabel.SetText(fmt.Sprintf("0x%04X", nandpu.RegJ.val))

		m1Label.SetText(fmt.Sprintf("0x%02X", nandpu.RegM.Lo.ForceGet()))
		m2Label.SetText(fmt.Sprintf("0x%02X", nandpu.RegM.Hi.ForceGet()))
		xLabel.SetText(fmt.Sprintf("0x%02X", nandpu.RegXY.Lo.ForceGet()))
		yLabel.SetText(fmt.Sprintf("0x%02X", nandpu.RegXY.Hi.ForceGet()))
		j1Label.SetText(fmt.Sprintf("0x%02X", nandpu.RegJ.Lo.ForceGet()))
		j2Label.SetText(fmt.Sprintf("0x%02X", nandpu.RegJ.Hi.ForceGet()))

		zeroLabel.SetText(fmt.Sprintf("%t", nandpu.Zero))
		carryLabel.SetText(fmt.Sprintf("%t", nandpu.Carry))
		signLabel.SetText(fmt.Sprintf("%t", nandpu.Sign))
		lessThanLabel.SetText(fmt.Sprintf("%t", nandpu.LessThan))

		stepNumLabel.SetText(fmt.Sprintf("Step: %d", stepNum))

		if running {
			runBtn.SetText("Stop")
			stepBtn.Disable()
			resetBtn.Disable()
		} else {
			if stepNum > 0 {
				resetBtn.Enable()
			} else {
				resetBtn.Disable()
			}
			runBtn.SetText("Run")
			stepBtn.Enable()
		}
	}

	go func() {
		for {
			if running {
				continueRunning := nandpu.Step()
				stepNum += 1
				fyne.Do(updateGUIValues)
				if !continueRunning {
					running = false
				}
				speedVal, err := speed.Get()
				if err != nil {
					Logger.Fatalf("Error: %s", err)
				}
				time.Sleep(time.Millisecond * time.Duration(speedVal))
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	Wnd.ShowAndRun()
}
