package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ImageInfo struct {
	Id  string
	URL string
}

func main() {
	http.HandleFunc("/", watchFlow)
	err := http.ListenAndServe(":21564", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func watchFlow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var t ImageInfo
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	downloadImage(t.Id, t.URL)
}

func downloadImage(fname string, url string) {
	err := os.MkdirAll("T:/DataStore/image_original", os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("T:/DataStore/image_original/%s.jpg", fname))

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[Image] Downloaded: %s\n", url)
}
