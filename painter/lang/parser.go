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
	lastBgRect  *painter.BgRectangle
	figures     []*painter.Figure
	moveOps     []painter.Operation
	updateOp    painter.Operation
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.init()
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation

	for scanner.Scan() {
		commandLine := scanner.Text()

		op, err := p.parse(commandLine)
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

func (p *Parser) parse(commandLine string) ([]painter.Operation, error) {
	parts := strings.Fields(commandLine)
	if len(parts) < 1 {
		return nil, nil
	}
	instruction := parts[0]
	iArgs, _ := getArgs(parts)

	switch instruction {
	case "white":
		p.lastBgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.lastBgColor = painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		if len(iArgs) != 4 {
			return nil, fmt.Errorf("invalid number of arguments for command bgrect: %v", len(iArgs))
		}
		p.lastBgRect = &painter.BgRectangle{X1: iArgs[0], Y1: iArgs[1], X2: iArgs[2], Y2: iArgs[3]}
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

func getArgs(parts []string) ([]int, error) {
	if len(parts) == 1 {
		return nil, nil
	}
	argsStr := parts[1:]
	var args []int
	for _, arg := range argsStr {
		i, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid argument: %v", arg)
		}
		args = append(args, i)
	}
	return args, nil
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
