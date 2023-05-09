package lang

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/EugeneSemivolos/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	lastBgColor painter.Operation
	lastBgRect  *painter.BgRect
	figures     []*painter.Figure
	moveOps     []painter.Operation
	updateOp    painter.Operation
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.init()
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()

		op, err := p.parseCommand(commandLine)
		if err != nil {
			return nil, err
		}
		res = append(res, op...)
	}
	return res, nil
}

func (p *Parser) init() {
	if p.lastBgColor == nil {
		p.lastBgColor = painter.OperationFunc(painter.ResetScreen)
	}
	if p.updateOp != nil {
		p.updateOp = nil
	}
}

func (p *Parser) parseCommand(commandLine string) ([]painter.Operation, error) {
	parts := strings.Fields(commandLine)
	if len(parts) < 1 {
		return nil, nil
	}
	instruction := parts[0]
	iArgs, err := parseArgs(parts[1:])
	if err != nil {
		return nil, err
	}

	switch instruction {
	case "white":
		p.lastBgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.lastBgColor = painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		if len(iArgs) != 4 {
			return nil, fmt.Errorf("invalid number of arguments for command bgrect: %v", len(iArgs))
		}
		p.lastBgRect = &painter.BgRect{X1: iArgs[0], Y1: iArgs[1], X2: iArgs[2], Y2: iArgs[3]}
	case "figure":
		if len(iArgs) != 2 {
			return nil, fmt.Errorf("invalid number of arguments for command figure: %v", len(iArgs))
		}
		clr := color.RGBA{R: 255, G: 255, B: 0, A: 1}
		figure := painter.Figure{X: iArgs[0], Y: iArgs[1], C: clr}
		p.figures = append(p.figures, &figure)
	case "move":
		if len(iArgs) != 2 {
			return nil, fmt.Errorf("invalid number of arguments for command move: %v", len(iArgs))
		}
		moveOp := painter.Move{X: iArgs[0], Y: iArgs[1], Figures: p.figures}
		p.moveOps = append(p.moveOps, &moveOp)
	case "reset":
		p.resetState()
		p.lastBgColor = painter.OperationFunc(painter.ResetScreen)
	case "update":
		p.updateOp = painter.UpdateOp
	default:
		return nil, fmt.Errorf("unknown command: %v", instruction)
	}

	return p.getParsedCommands(), nil
}

func parseArgs(argsStr []string) ([]float32, error) {
	if len(argsStr) == 0 {
		return nil, nil
	}
	var res []float32

	for _, argStr := range argsStr {
		arg, err := strconv.ParseFloat(argStr, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid params")
		}

		if arg < 0 || arg > 1 {
			return nil, fmt.Errorf("invalid coordinates")
		}

		res = append(res, float32(arg))
	}

	return res, nil
}

func (p *Parser) resetState() {
	p.lastBgColor = nil
	p.lastBgRect = nil
	p.figures = nil
	p.moveOps = nil
	p.updateOp = nil
}

func (p *Parser) getParsedCommands() []painter.Operation {
	var res []painter.Operation
	if p.lastBgColor != nil {
		res = append(res, p.lastBgColor)
	}
	if p.lastBgRect != nil {
		res = append(res, p.lastBgRect)
	}
	if len(p.moveOps) != 0 {
		res = append(res, p.moveOps...)
	}
	p.moveOps = nil
	if len(p.figures) != 0 {
		for _, figure := range p.figures {
			res = append(res, figure)
		}
	}
	if p.updateOp != nil {
		res = append(res, p.updateOp)
	}
	return res
}
