package main

import (
	"bufio"
	"fmt"
	"go_lsp/rpc"
	"log"
	"os"
)

func main() {
	fmt.Print("Hi !")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(rpc.Split)
	logger := getLogger("/home/alvi/Desktop/log.txt")
	logger.Println("Hello World !")

	for scanner.Scan() {
		message := scanner.Text()
		handleMessage(logger, message)

	}
}

func handleMessage(logger *log.Logger, message any) {
	logger.Println(message)
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didnt give me a good file")
	}
	return log.New(logfile, "[go_lsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
