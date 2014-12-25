package main

import (
	"net/http"
	"encoding/json"
	"os"
	"log"
	"path/filepath"
	"regexp"
)

type Configuration struct {
	Port    string
}

func main() {
	file, err := os.Open("server.json")
	configuration := Configuration{Port: "7654"}

	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&configuration)

		if err != nil {
			log.Fatal("Could not decode config: ", err)
			return
		}
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	p := ":" + configuration.Port

	log.Printf("Listening on port %s.", configuration.Port)
	log.Printf("Open http://localhost%s/ in your browser.", p)
	log.Printf("Current working directory: %s", dir)

	os.Chdir(dir)

	log.Fatal(http.ListenAndServe(p, handleFileServer(dir)))
}

func handleFileServer(dir string) http.HandlerFunc {
	realHandler := http.FileServer(http.Dir(dir)).ServeHTTP

	fileRegexp := regexp.MustCompile(`^(.+)\.(.+)$`)

	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(dir + r.URL.Path); os.IsNotExist(err) {
			log.Println("File not found:", r.URL.Path)
			if !fileRegexp.MatchString(r.URL.Path) {
				notFound(w, r)
				return
			}
		}
		realHandler(w, r)
	}
}


func notFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#!"+string(r.URL.Path), http.StatusSeeOther)
}
