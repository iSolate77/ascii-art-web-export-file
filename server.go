package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ASCIIArt struct {
	Text   string
	Banner string
	Result string
}

var ascii_art_to_download string

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, ASCIIArt{})
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if len(text) > 100 {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	if text == "" || banner == "" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	result := generateASCIIArt(text, banner)
	ascii_art_to_download = result

	if result == "" {
		if banner == "internalServerError" {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	asciiArt := ASCIIArt{Text: text, Banner: banner, Result: result}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, asciiArt)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusInternalServerError {
		tmpl, _ := template.ParseFiles("./templates/500.html")
		tmpl.Execute(w, nil)
	} else if status == http.StatusNotFound {
		tmpl, _ := template.ParseFiles("./templates/404.html")
		tmpl.Execute(w, nil)
	} else if status == http.StatusBadRequest {
		tmpl, _ := template.ParseFiles("./templates/400.html")
		tmpl.Execute(w, nil)
	}
}

func generateASCIIArt(text string, banner string) string {
	var standard_file string
	if banner == "standard" {
		standard_file = "./banners/standard.txt"
	} else if banner == "shadow" {
		standard_file = "./banners/shadow.txt"
	} else if banner == "thinkertoy" {
		standard_file = "./banners/thinkertoy.txt"
	} else if banner == "internalServerError" {
		standard_file = "./banners/internalServerError.txt"
	} else {
		fmt.Println("Banner not found")
		return ""
	}

	ascii_art := []string{}
	file_reader, err := os.Open(standard_file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file_reader.Close()

	scanner := bufio.NewScanner(file_reader)
	for scanner.Scan() {
		ascii_art = append(ascii_art, scanner.Text())
	}

	var output []string
	for _, part := range strings.Split(text, "\n") {
		if part == "" {
			output = append(output, "\n")
			continue
		}

		ascii_values := []int{}
		for _, c := range part {
			ascii_values = append(ascii_values, int(c))
		}

		for j := 0; j < 8; j++ {
			for _, val := range ascii_values {
				if val != 10 && val != 13 && (val < 32 || val > 126) {
					fmt.Println("ASCII value out of range for:", val)
					return ""
				}

				line_number := 9*val - 287
				if line_number < 0 || line_number >= len(ascii_art) {
					// handle cases where ASCII value is out of range
					continue
				}
				output = append(output, ascii_art[line_number+j])
			}
			output = append(output, "\n")
		}
	}

	return strings.Join(output, "")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")
	w.Header().Set("Content-Length", strconv.Itoa(len(ascii_art_to_download)))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(ascii_art_to_download))
}

func main() {
	fs := http.FileServer(http.Dir("./templates/"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			mainPageHandler(w, r)
		case "/ascii-art":
			asciiArtHandler(w, r)
		case "/download-ascii":
			downloadHandler(w, r)
		default:
			errorHandler(w, r, http.StatusNotFound)
		}
	})

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
