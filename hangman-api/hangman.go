package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var inputletter []string

/*Prend la lettre ou le mot choisi et vérifie si le mot est complet sinon renvoie le dessin
correspondant etle nombre de chance restant */
func hangman(hgn hangm, r *http.Request) {
	textascii := []string{}
	file, err := os.Open("standard.txt")
	if err != nil {
		log.Fatalf("failed to open standard.txt")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		textascii = append(textascii, scanner.Text()) // Slice ascii
	}
	strsolution := ""
	for _, v := range hgn.Solution { //Incrémentation de la solutin
		strsolution += string(v)
	}
	if hgn.Chance == 10 {
		inputletter = []string{}
		//1-----------------fmt.Println("Good Luck, you have", hgn.Chance, "attempts.")
		asciiletter(textascii, hgn)
	} else {
		//1-----------------fmt.Println("You have", hgn.Chance, "attempts remaining.")
	}
	equal := true
	if hgn.Chance > 0 {
		hgn.Doublon = false
		for i := 0; i < len(hgn.Word); i++ { //vérifie si le mot est complet
			if hgn.Word[i] != hgn.Solution[i] {
				equal = false
				break
			}
		}
		if equal == true { //si le mot est complet victoire et arrêt
			//1-----------------fmt.Println("Congrats !")
			hgn.Chance = 0
			asciiletter(textascii, hgn)
			dessin(hgn)
			majstruct(hgn)
			os.Remove("resultat.txt")
			return
		}
		//1-----------------fmt.Print("\n\nChoose: ")
		//-----------------------------------------------------------------	add param to URL "localhost:8080/hangman?key=a"
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			fmt.Print("le chemin", r.URL)
			return
		}
		key := keys[0]
		input := key
		//-----------------------------------------------------------------
		here := false
		if len(input) > 1 && input != "STOP" && hgn.Chance > 1 { //Traite input == mot
			if input == strsolution {
				//1-----------------fmt.Println("Congrats !")
				asciiletter(textascii, hgn)
				dessin(hgn)
				majstruct(hgn)
				os.Remove("resultat.txt")
				return
			} else {
				hgn.Chance -= 2
				if hgn.Chance == 8 {
					hgn.Drawingpos += 8
				} else {
					hgn.Drawingpos += 16
				}
				//1-----------------fmt.Println("\nNot present in the word, ", hgn.Chance, " attempts remaining")
			}
		} else { //Traite input == lettre
			inputletter = append(inputletter, input)
			for i := range inputletter {
				if i != 0 && inputletter[i-1] == input {
					//1-----------------fmt.Println("Error letter already submitted")
					hgn.Doublon = true
				} else {
					for i, v := range hgn.Solution {
						if input == v {
							//1-----------------fmt.Println(string(v))
							hgn.Word[i] = string(v)
							here = true
						}
					}
				}
			}
		}
		if !here && len(input) <= 1 && !hgn.Doublon { //pas présent print chance-1 et dessin correspondant
			hgn.Chance--
			//1-----------------fmt.Println("Not present in the word, ", hgn.Chance, " attempts remaining")
			fmt.Println("- Bad letter")
			if hgn.Chance != 9 {
				hgn.Drawingpos += 8
			}
		}
		hgn.Ascii = asciiletter(textascii, hgn)
		hgn.Drawing = dessin(hgn)
		majstruct(hgn)
	}
	if hgn.Chance == 0 && !equal {
		//1-----------------fmt.Println("You lost, try again")
	}
	os.Remove("resultat.txt")
}

func majstruct(hgn hangm) {
	structure, _ := json.Marshal(hgn)
	f, _ := os.Create("Struct.json")
	f.Write(structure)
}

//Prend le mot en parametre et renvoie le résultat en lettres acscii passant en paramètres
func asciiletter(ascii []string, hgn hangm) []string {
	file, err := os.OpenFile("resultat.txt", os.O_CREATE|os.O_WRONLY, 0600)
	defer file.Close() // on ferme automatiquement à la fin de notre programme
	//L'objectif ici est de donner une valeur (un dessin) à la letrre "OS.arg[:1]"
	slicealpha := []string{}
	alpha := "abcdefghijklmnopqrstuvwxyz"
	count := 0
	essaie := []string{}
	//Ici min et max correspondent aux index séparant les lettre ascii que nous imprimerons par la suite
	for _, x := range alpha {
		slicealpha = append(slicealpha, string(x))
	}
	for i := range hgn.Word {
		/*
			Si la lettre recu correspond à une dès lettre de lalphabet contenue dans le slice (slicealpha)
			alors (count) calcul la position de la lettre en question
		*/
		for j := range slicealpha {
			count++

			if hgn.Word[i] == slicealpha[j] {
				min := (count * 9) + 577
				max := (count * 9) + 586

				for _, x := range ascii[min:max] {
					_, err = file.WriteString("\n")
					_, err = file.WriteString(string(x))
					// écrire dans le fichier "resultat.txt" la lettre correspondante
					if err != nil {
						panic(err)
					}
				}
			}
		}
		if hgn.Word[i] != hgn.Solution[i] {
			for _, x := range ascii[116:125] {
				_, err = file.WriteString("\n")
				_, err = file.WriteString(string(x))
				// écrire dans le fichier "resultat.txt" le caractère "_"
				if err != nil {
					panic(err)
				}
			}

		}
		count = 0
	}
	counter := 9
	counter2 := 1
	tabascii := []string{}
	for j := 0; j < 9; j++ { //Recupère l'index+1 \n de chacune des lettres dans "resultat.txt"
		for i := len(hgn.Solution); i > 0; i-- {
			counter2++
			n := 0
			n = 9*counter2 - counter
			filer, _ := os.Open("resultat.txt")
			scanner := bufio.NewScanner(filer)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				essaie = append(essaie, scanner.Text())
			}
			//1-----------------fmt.Print(essaie[n])
			tabascii = append(tabascii, essaie[n])
		}
		if j != 9 {
			//1-----------------fmt.Print("\n")
		}
		counter2 = 0
		counter--
	}
	return tabascii
}

//print dessin correspondant
func dessin(hgn hangm) []string {
	var text []string
	file, err := os.Open("hangman.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file) //lis le fichier qui contient les dessins
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	tabdraw := []string{}
	for _, x := range text[hgn.Drawingpos : hgn.Drawingpos+8] { //Print le dessin ,via la dernière position drawingpos(l'index) défini dans la struct +8
		//1-----------------fmt.Println(string(x))
		tabdraw = append(tabdraw, string(x)+"\n") //stocker pr renvoyer en web le dessin dans tabdraw
	}
	return tabdraw
}

//initialisation du mot incomplet et de sa solution
func printn(word string, mi int, ma int) ([]string, []string) {
	wordarr := make([]string, len(word))
	ma = 0
	solution := []string{}
	for _, v := range word { //met le mot complet dans solution
		solution = append(solution, string(v))
		ma++
	}
	for i := 0; i < len(word)/2-1; i++ { //Révèle n lettres random du mot où n est len(word) / 2 - 1
		r := rand.Intn(ma-mi-1) + mi
		if wordarr[r] == "" {
			wordarr[r] = solution[r]
		}
	}

	for i := 0; i < len(wordarr); i++ { //Remplace chaque case vide par un "_"
		if wordarr[i] == "" {
			wordarr[i] = "_"
		}
	}
	return wordarr, solution
}

func hangmain(w http.ResponseWriter, r *http.Request, hgn hangm) {
	save2, _ := ioutil.ReadFile("Struct.json") //reprend sauvegarde a chaque appel
	err := json.Unmarshal([]byte(save2), &hgn)
	if err != nil {
		log.Fatalf("failed to encode game")
	}

	keys := r.URL.Query()["key"] //prend la key renseignée dans l'url pour lancer avec niveau voulu
	inpute := keys[0]
	if hgn.Chance <= 0 || inpute == "Level1" || inpute == "Level2" || inpute == "Level3" || inpute == "Retry" {
		var text []string
		min := 0
		max := 0
		fmt.Println("New game start")
		//-----------------------------------------------------------------
		if inpute == "Level1" { //DEFINIR NIVEAU
			fmt.Println("- Level 1 selected")
			file, _ := os.Open("words.txt")
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				text = append(text, scanner.Text())
				max++ //borne max random
			}
		} else if inpute == "Level2" { //DEFINIR NIVEAU
			fmt.Println("- Level 2 selected")
			file, _ := os.Open("words2.txt")
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				text = append(text, scanner.Text())
				max++ //borne max random
			}
		} else {
			fmt.Println("- Level 3 selected")
			file, _ := os.Open("words3.txt") //DEFINIR NIVEAU
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				text = append(text, scanner.Text())
				max++ //borne max random
			}
		}
		if inpute == "- Retry" {
			fmt.Println("retry selected")
			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(max-min) + min
			word := text[random]
			tmpword, tmpsolution := (printn(word, min, max))
			hgn := hangm{tmpsolution, tmpword, 12, 0, nil, nil, false}
			hangman(hgn, r)
		} else {
			fmt.Println("- reset selected")
			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(max-min) + min
			word := text[random]
			tmpword, tmpsolution := (printn(word, min, max))
			hgn := hangm{tmpsolution, tmpword, 12, 0, nil, nil, false}
			hangman(hgn, r)
		}
	} else {
		hangman(hgn, r)
	}
}
