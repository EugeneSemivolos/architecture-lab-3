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

type BgRectangle struct {
	X1, Y1, X2, Y2 int
}

func (op *BgRectangle) Do(tx screen.Texture) bool {
	tx.Fill(image.Rect(op.X1, op.Y1, op.X2, op.Y2), color.Black, screen.Src)
	return false
}

type Figure struct {
	X, Y int
	C    color.RGBA
}

func (op *Figure) Do(tx screen.Texture) bool {
	tx.Fill(image.Rect(op.X-150, op.Y-50, op.X+150, op.Y+50), op.C, draw.Src)
	tx.Fill(image.Rect(op.X-50, op.Y-150, op.X+50, op.Y+150), op.C, draw.Src)
	return false
}

type Move struct {
	X, Y    int
	Figures []*Figure
}

func (op *Move) Do(tx screen.Texture) bool {
	for i := range op.Figures {
		op.Figures[i].X += op.X
		op.Figures[i].Y += op.Y
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
