package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		fmt.Println(err)
		return
	}

	fmt.Println("Listening: ", 21564)
}

func loadConfig() {
	filePtr, err := os.Open("slings-config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer filePtr.Close()
	var info Config
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		fmt.Println("Failed loading config: ", err.Error())
		return
	}

	fmt.Println("DataStore: ", info.DataStore)

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
		fmt.Println(err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.jpg", imgpath, fname))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("[Image] Downloaded: %s\n", url)
}
