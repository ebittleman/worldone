package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func Add(a int, b int) int {
	return a + b
}

func main() {
	// var i int
	// _, err := fmt.Scanf("%d", &i)
	// fmt.Print(err)

	logger := log.New(os.Stderr, "helloword", log.Ldate|log.Ltime)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Fatal(err)
			continue
		}
		com := &Comm{conn, make(chan io.ReadWriteCloser)}
		go handleConnection(com)
		com.out <- os.Stderr
	}
}

type Comm struct {
	in  io.ReadWriteCloser
	out chan io.ReadWriteCloser
}

func handleConnection(com *Comm) {
	done := make(chan bool)

	output := <-com.out
	input := com.in

	sendFunc := func(input io.ReadWriteCloser, output io.Writer) {

		defer input.Close()
		reader := bufio.NewReader(input)

		defer func() {
			done <- true
		}()

		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = output.Write([]byte(line))
			if err != nil {
				fmt.Println(err)
				return
			}
			input.Write([]byte("thanks\n"))
		}
	}

	go sendFunc(input, output)
	<-done
	close(done)
	fmt.Println("All Closed Up Later")
}
