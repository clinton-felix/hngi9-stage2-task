package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Declaring the input data..
type InputData struct {
	OperationType string `json:"operation_type"`
	Num1          int64  `json:"x"`
	Num2          int64  `json:"y"`
}

// Declaring the output data struc format
type OutputData struct {
	SlackUsername	string		`json:"slackUsername"`
	Result			int64		`json:"result"`
	OperationType	string	`json:"operation_type"`
}

func scan(myInput InputData) (int64, string) {
	var operation int64
	var operation_type string

	operands := map[string][]string{
		"addition": {"addition", "add", "plus", "sum"},
		"subtraction": {"subtract", "subtraction", "minus"},
		"multiplication": {"multiply", "multiplication", "times", "product"},
	}

	for _, val := range operands["addition"] {
		if strings.Contains(myInput.OperationType, val) {
			operation = myInput.Num1 + myInput.Num2
			operation_type = "addition"
		}
	}
	
	for _, val := range operands["subtraction"] {
		if strings.Contains(myInput.OperationType, val) {
			operation = myInput.Num1 - myInput.Num2
			operation_type = "subtraction"
		}
	}

	for _, val := range operands["multiplication"] {
		if strings.Contains(myInput.OperationType, val) {
			operation = myInput.Num1 * myInput.Num2
			operation_type = "multiplication"
		}
	}

	return operation, operation_type
}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
 }

func OperationFunc(w http.ResponseWriter, r *http.Request){
	// Handling CORS
	setupCorsResponse(&w, r)
		if (*r).Method == "OPTIONS" {
		   return
		}

	var myInput InputData
	
	// encoding the request
	err := json.NewDecoder(r.Body).Decode(&myInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, operation_type := scan(myInput)

	output := OutputData{
		SlackUsername: "godDev",
		Result: result,
		OperationType: operation_type,
	}

	// marshalling the output into json
	outputByte, _ := json.Marshal(output)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(outputByte)
}


func main() {
	c := http.NewServeMux()

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading the .env File...")
	}
	port := os.Getenv("PORT")
	if port == ""{
		port = "3000"
	}

	
	c.HandleFunc("/", OperationFunc)
	
	// Running the server..
	fmt.Printf("Listening on port %v...", port)
	log.Fatal(http.ListenAndServe(":"+port, c))
	http.Handle("/", c)
}