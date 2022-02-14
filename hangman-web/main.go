package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//login/multijoueur
//score /scoreboard

var templatesDir = os.Getenv("TEMPLATES_DIR")

type hangm struct {
	Solution   []string
	Word       []string
	Chance     int
	Drawingpos int
	Drawing    []string
	Ascii      []string
	Doublon    bool
}

type Html struct {
	Arrstr template.HTML //ligne html dessin pendu
}
type Finalstruct struct {
	Arrstr      template.HTML
	Solution    template.HTML
	Word        template.HTML
	Chance      int
	Drawingpos  int
	Drawing     []string
	Ascii       []string
	Doublon     bool
	Solutionstr string
	Wordstr     string
}

var stop int = 0

func send(w http.ResponseWriter, r *http.Request) {
	//-------------------------- Récup param
	if keys, ok := r.URL.Query()["key"]; !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	} else if keys[0] == "Reset" {
		http.Redirect(w, r, "/hangman_welcome.html", http.StatusFound) //Return sur la page initial
	} else {
		fmt.Println("New param sent to API : http://localhost:8080/hangmanPOST?key=" + keys[0])
		//-------------------------- Fin récup param
		//-------------------------- Début envoi à l'api
		url := "http://localhost:8080/hangmanPOST?key=" + keys[0]
		resp, err := http.Post(url, keys[0], nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		//-------------------------- Fin envoie à l'api
	}
	//-------------------------- Début de récup json
	zero := []string{}
	hgn := hangm{zero, zero, 11, 0, nil, nil, false} //creation instance struct pr stocker resultat requete
	url := "http://localhost:8080/hangmanGET"
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "spacecount-tutorial")
	res, _ := spaceClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &hgn)
	//-------------------------- Fin récup json
	//-------------------------- Dessin Html (récup)//--------------------------  convertir en template.HTML pour envoyer tab sans []
	ht := Html{""}
	ht.Arrstr = draw(ht, hgn)
	strword := ""
	strsolution := ""
	for _, letter := range hgn.Solution {
		strsolution += letter
	}
	for _, letter := range hgn.Word {
		strword += letter
	}
	//-------------------------- Fin returns string(dessin)
	//--------------------------  convertir en template.HTML pour envoyer tab sans []
	sol := ""
	wor := ""
	for i := 0; i < len(hgn.Word); i++ {
		wor += hgn.Word[i]
	}
	for i := 0; i < len(hgn.Solution); i++ {
		sol += hgn.Solution[i]
	}
	wor2 := template.HTML(wor)
	sol2 := template.HTML(sol)
	//--------------------------  Fin convertir en template.HTML pour envoyer tab sans []
	final := Finalstruct{ht.Arrstr, sol2, wor2, hgn.Chance, hgn.Drawingpos, hgn.Drawing, hgn.Ascii, hgn.Doublon, strsolution, strword} //struct a envoyer au site
	stop++
	for i := range hgn.Word {
		if hgn.Word[i] != hgn.Solution[i] {
			stop = 0
			break
		}
	}
	//-------------------------- Fin cas de victoire
	//-------------------------- Debut injection template
	if stop == 2 { //win
		http.Redirect(w, r, "/hangman_welcome.html", http.StatusFound)
	}
	if final.Chance < 0 { //lose
		http.Redirect(w, r, "/hangman_welcome.html", http.StatusFound) //Return sur la page initial
	}
	tmplPath := filepath.Join(templatesDir, "static/hangman.html")
	tmpl := template.Must(template.ParseFiles(tmplPath))
	tmpl.Execute(w, final)
	//-------------------------- Fin injection template
}

func draw(ht Html, hgn hangm) template.HTML {
	//-------------------------- Récup dessin (problème <>)
	file, _ := os.Open("static/bringdeath.txt")
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	arr := []string{}
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	//-------------------------- Fin récup dessin
	//-------------------------- Début récup bon dessin selon nb chances
	count := hgn.Chance - 10
	if count < 0 {
		count *= -1
	}
	array := ""
	for i := 0; i < count; i++ {
		array += arr[i] + "\n"
	}
	t := template.HTML(array)
	return t
	//-------------------------- Fin récup bon dessin et return
}

func main() {
	fs := http.FileServer(http.Dir("./staticc/"))
	http.Handle("/", fs)
	http.HandleFunc("/hangmanPOST", send)
	fmt.Printf("Starting server at port 8189 successfully\n")
	http.ListenAndServe(":8189", nil)
}
