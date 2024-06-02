package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type client chan<- string

var (
	entering   = make(chan client)
	leaving    = make(chan client)
	expression = make(chan string)
	messages   = make(chan string)
	answer     = make(chan int)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}
	go broadcaster()
	go calculations()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go connHandler(conn)
	}
}

func calculations() {
	var (
		operation string
		ans       int
	)

	for {
		operand1 := rand.Intn(1000)
		operand2 := rand.Intn(1000)
		n := rand.Intn(4)

		switch n {
		case 0:
			operation = " + "
			ans = operand1 + operand2
		case 1:
			operation = " - "
			ans = operand1 - operand2
		case 2:
			operation = " * "
			ans = operand1 * operand2
		case 3:
			operation = " / "
			if operand1 < operand2 {
				operand1, operand2 = operand2, operand1
			}
			for {
				if operand2 == 0 {
					operand2 = rand.Intn(operand1)
				} else {
					break
				}
			}
			ans = operand1 / operand2
		}
		exp := fmt.Sprint(operand1) + operation + fmt.Sprint(operand2)
		answer <- ans
		expression <- exp
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case exp := <-expression:
			for cli := range clients {
				cli <- exp
			}
		case ans := <-answer:

		}
	}
}

func connHandler(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Successfully connected\nEnter your nickname:\n"))
	var (
		ch   = make(chan string)
		buf  = make([]byte, 256)
		flag bool
	)
	go clientWriter(conn, ch)
	_, err := conn.Read(buf)
	if err != nil {
		log.Print(err)
	}
	who := string(buf)
	who = strings.Replace(who, "\n", "", -1)
	ch <- "You are " + who
	messages <- who + " has entered"
	entering <- ch
	input := bufio.NewScanner(conn)
	for {
		rightAnswer := <-answer
		for {
			if flag = input.Scan(); !flag {
				break
			}
			clientAnswer, err := strconv.Atoi(input.Text())
			if err != nil {
				ch <- "Invalid format!"
				continue
			}
			messages <- who + ": " + input.Text()
			if rightAnswer != clientAnswer {
				messages <- "Wrong answer!"
			} else {
				messages <- who + " win!"
				break
			}
		}
		if !flag {
			break
		}
	}
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
