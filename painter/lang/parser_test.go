package lang

import (
	"image/color"
	"strings"
	"testing"

	"github.com/EugeneSemivolos/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Test struct {
	name    string
	command string
	op      painter.Operation
}

var testsStruct = []Test{
	{
		name:    "background rectangle",
		command: "bgrect 0 0 0.1 0.1",
		op:      &painter.BgRect{X1: 0, Y1: 0, X2: 0.1, Y2: 0.1},
	},
	{
		name:    "figure",
		command: "figure 0.25 0.25",
		op:      &painter.Figure{X: 0.25, Y: 0.25, C: color.RGBA{R: 255, G: 255, B: 0, A: 1}},
	},
	{
		name:    "move",
		command: "move 0.3 0.4",
		op:      &painter.Move{X: 0.3, Y: 0.4},
	},
	{
		name:    "update",
		command: "update",
		op:      painter.UpdateOp,
	},
	{
		name:    "invalid command",
		command: "invalidcommand",
		op:      nil,
	},
}

var testsFunc = []Test{
	{
		name:    "white fill",
		command: "white",
		op:      painter.OperationFunc(painter.WhiteFill),
	},
	{
		name:    "green fill",
		command: "green",
		op:      painter.OperationFunc(painter.GreenFill),
	},
	{
		name:    "reset screen",
		command: "reset",
		op:      painter.OperationFunc(painter.ResetScreen),
	},
}

func TestParseStruct(t *testing.T) {
	for _, tc := range testsStruct {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))
			if tc.op == nil {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.IsType(t, tc.op, ops[1])
				assert.Equal(t, tc.op, ops[1])
			}
		})
	}
}

func TestParseFunc(t *testing.T) {
	parser := &Parser{}

	for _, tc := range testsFunc {
		t.Run(tc.name, func(t *testing.T) {
			ops, err := parser.Parse(strings.NewReader(tc.command))
			require.NoError(t, err)
			require.Len(t, ops, 1)
			assert.IsType(t, tc.op, ops[0])
		})
	}
}
