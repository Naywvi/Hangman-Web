package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

type hangm struct {
	Solution   []string //mot solution
	Word       []string //mot incomplet
	Chance     int      //nb chance
	Drawingpos int      //dessin correspondant au nb de chance(s)
	Drawing    []string
	Ascii      []string
	Doublon    bool //pour gerer les doublons en web
}

func hangmanOutputHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hangmanGET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusBadRequest)
		return
	} else {
		fmt.Println("New GET request")
		//x--------------------------------------------------x Open Json handler
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		//x--------------------------------------------------x
		aff, err := os.Open("Struct.json")
		if err != nil {
			log.Fatalf("failed to encode Struct.json in httpServ")
		}
		scan := bufio.NewScanner(aff)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			w.Write([]byte(scan.Text())) //xPrint byte
		}
	}

}
func hangmanstartHandler(w http.ResponseWriter, r *http.Request) { // "localhost:8080/hangman"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("New POST request : ", r.URL)
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusBadRequest)
		return
	} else if r.Method == "POST" {
		hangmain(w, r, hangm{})
	}
}

func main() {
	FileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", FileServer)
	http.HandleFunc("/hangmanGET", hangmanOutputHandler)
	http.HandleFunc("/hangmanPOST", hangmanstartHandler)

	fmt.Printf("Starting server at port 8080 successfully\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
