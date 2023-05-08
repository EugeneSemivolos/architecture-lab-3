package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	wnd  screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz       size.Event
	pos      image.Rectangle
	mousePos image.Point
}

func (pv *Visualizer) Main() {
	pv.tx = make(chan screen.Texture)
	pv.done = make(chan struct{})
	pv.pos.Max.X = 200
	pv.pos.Max.Y = 200
	pv.mousePos = image.Point{X: 400, Y: 400}
	driver.Main(pv.run)
}

func (pv *Visualizer) Update(tx screen.Texture) {
	pv.tx <- tx
}

func (pv *Visualizer) run(scr screen.Screen) {
	if pv.OnScreenReady != nil {
		pv.OnScreenReady(scr)
	}

	wnd, err := scr.NewWindow(&screen.NewWindowOptions{
		Title:  pv.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		wnd.Release()
		close(pv.done)
	}()

	pv.wnd = wnd

	events := make(chan any)
	go func() {
		for {
			e := wnd.NextEvent()
			if pv.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var tx screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pv.handleEvent(e, tx)

		case tx = <-pv.tx:
			wnd.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pv *Visualizer) handleEvent(e any, tx screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pv.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if tx == nil {
			// Реалізація реакції на натискання кнопки миші.
			if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
				pv.mousePos = image.Point{
					X: int(e.X),
					Y: int(e.Y),
				}
				pv.wnd.Send(paint.Event{})
				fmt.Printf("Clicked on point: x=%d, y=%d\n", int(e.X), int(e.Y))
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if tx == nil {
			pv.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pv.wnd.Scale(pv.sz.Bounds(), tx, tx.Bounds(), draw.Src, nil)
		}
		pv.wnd.Publish()
	}
}

func (pv *Visualizer) drawDefaultUI() {
	X, Y := pv.mousePos.X, pv.mousePos.Y
	fillColor := color.RGBA{R: 255, G: 255, B: 0, A: 1}

	pv.wnd.Fill(pv.sz.Bounds(), color.Black, draw.Src) // Фон.

	// TODO: Змінити колір фону та додати відображення фігури у вашому варіанті.
	pv.wnd.Fill(image.Rect(X-150, Y-50, X+150, Y+50), fillColor, draw.Src)
	pv.wnd.Fill(image.Rect(X-50, Y-150, X+50, Y+150), fillColor, draw.Src)
	// Малювання білої рамки.
	for _, br := range imageutil.Border(pv.sz.Bounds(), 10) {
		pv.wnd.Fill(br, color.White, draw.Src)
	}
}
