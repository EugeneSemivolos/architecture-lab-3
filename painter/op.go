package painter

import (
	"image/color"

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

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(tx screen.Texture) {
	tx.Fill(tx.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(tx screen.Texture) {
	tx.Fill(tx.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}
