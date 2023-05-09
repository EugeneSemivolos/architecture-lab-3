package painter

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(tx screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (opList OperationList) Do(tx screen.Texture) (ready bool) {
	for _, op := range opList {
		ready = op.Do(tx) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(tx screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(tx screen.Texture) bool {
	f(tx)
	return false
}

type BgRect struct {
	X1, Y1, X2, Y2 float32
}

func (op *BgRect) Do(tx screen.Texture) bool {
	tx.Fill(
		image.Rect(
			int(op.X1*float32(tx.Size().X)),
			int(op.Y1*float32(tx.Size().Y)),
			int(op.X2*float32(tx.Size().X)),
			int(op.Y2*float32(tx.Size().Y)),
		),
		color.Black,
		screen.Src,
	)
	return false
}

type Figure struct {
	X, Y float32
	C    color.RGBA
}

func (op *Figure) Do(tx screen.Texture) bool {
	X0 := int(op.X * float32(tx.Size().X))
	Y0 := int(op.Y * float32(tx.Size().Y))
	tx.Fill(
		image.Rect(X0-150, Y0-50, X0+150, Y0+50),
		op.C,
		draw.Src,
	)
	tx.Fill(
		image.Rect(X0-50, Y0-150, X0+50, Y0+150),
		op.C,
		draw.Src,
	)
	return false
}

type Move struct {
	X, Y    float32
	Figures []*Figure
}

func (op *Move) Do(tx screen.Texture) bool {
	for i := range op.Figures {
		op.Figures[i].X = op.X
		op.Figures[i].Y = op.Y
	}
	return false
}

// ResetScreen зафарбовує тестуру у чорний колір.
func ResetScreen(tx screen.Texture) {
	tx.Fill(tx.Bounds(), color.Black, draw.Src)
}

// WhiteFill зафарбовує тестуру у білий колір.
func WhiteFill(tx screen.Texture) {
	tx.Fill(tx.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір.
func GreenFill(tx screen.Texture) {
	tx.Fill(tx.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}
