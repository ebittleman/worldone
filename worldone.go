package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"code.google.com/p/go.net/websocket"
)

func Add(a int, b int) int {
	return a + b
}

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s *socket) Close() error {
	close(s.done)
	return nil
}

var listenAddr = "localhost:8080"

func main() {
	// var i int
	// _, err := fmt.Scanf("%d", &i)
	// fmt.Print(err)

	logger := log.New(os.Stderr, "helloword", log.Ldate|log.Ltime)

	go telnetServer(logger)

	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		logger.Fatal(err)
	}
}

func telnetServer(logger *log.Logger) {
	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Fatal(err)
			continue
		}
		go match(conn)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<script>
function onMessage(m){
	div = document.createElement("div");
	div.innerHTML = m.data;
	input.appendChild(div)
}
function onClose(){
	div = document.createElement("div");
	div.innerHTML = "BYE!";
	input.appendChild(div)
}
websocket = new WebSocket("ws://{{.}}/socket");
websocket.onmessage = onMessage;
websocket.onclose = onClose;
</script>
<body>
<div id="input"></div><br/>
<input type="text" id="ouput"/><button id="btn">Send</button>
</body>
</html>    
</html>
`))

func socketHandler(ws *websocket.Conn) {
	s := &socket{ws, make(chan bool)}
	go match(s)
	<-s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	logger := log.New(os.Stderr, "helloword", log.Ldate|log.Ltime)
	fmt.Fprint(c, "Waiting ...")
	select {
	case partner <- c:
	case p := <-partner:
		chat(p, c, logger)
	}
}

func chat(a, b io.ReadWriteCloser, logger *log.Logger) {
	fmt.Fprintln(a, "Welcome")
	fmt.Fprintln(b, "Welcome")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)

	if err := <-errc; err != nil {
		logger.Println(err)
	}

	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
