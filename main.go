package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go_lsp/analysis"
	"go_lsp/lsp"
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
	// new state
	state := analysis.NewState()

	for scanner.Scan() {
		message := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(message)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, method, contents, state)

	}
}

func handleMessage(logger *log.Logger, method string, contents []byte, state analysis.State) {
	logger.Printf("Received message with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Parser can not parse this: %s", err)
		}
		logger.Printf("Connected to %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)
		// Let's reply
		message := lsp.NewInitializedResponse(request.ID)
		reply := rpc.EncodeMessage(message)
		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Sent the reply")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, could not parse this %s", err)
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, could not parse this %s", err)
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didnt give me a good file")
	}
	return log.New(logfile, "[go_lsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
