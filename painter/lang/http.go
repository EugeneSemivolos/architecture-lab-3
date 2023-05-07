package lang

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/EugeneSemivolos/architecture-lab-3/painter"
)

// HttpHandler конструює обробник HTTP запитів, який дані з запиту віддає у Parser, а потім відправляє отриманий список
// операцій у painter.Loop.
func HttpHandler(loop *painter.Loop, parser *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var in io.Reader = req.Body
		if req.Method == http.MethodGet {
			in = strings.NewReader(req.URL.Query().Get("cmd"))
		}

		cmds, err := parser.Parse(in)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		loop.Post(painter.OperationList(cmds))
		rw.WriteHeader(http.StatusOK)
	})
}
