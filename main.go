package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	DataStore string
}

type ImageInfo struct {
	Id  string
	URL string
}

var dpath string

func main() {
	loadConfig()

	http.HandleFunc("/", watchFlow)
	err := http.ListenAndServe(":21564", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func loadConfig() {
	filePtr, err := os.Open("slings-config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer filePtr.Close()
	var info Config
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		fmt.Println("Failed loading config: ", err.Error())
	}

	dpath = info.DataStore
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
	imgpath := fmt.Sprintf("%s/image_original", dpath)

	err := os.MkdirAll(imgpath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.jpg", imgpath, fname))

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
