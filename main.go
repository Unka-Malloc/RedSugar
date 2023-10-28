package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Config struct {
	DataStore string `json:"DataStore"`
}

type ImageInfo struct {
	Id  string `json:"img_id"`
	URL string `json:"img_url"`
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
		fmt.Println("[Err] Loading Config: ", err)
		return
	}

	fmt.Println("DataStore: ", info.DataStore)

	dpath = info.DataStore
}

func watchFlow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var imgArr []ImageInfo
	err := decoder.Decode(&imgArr)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(imgArr); i++ {
		go downloadImage(imgArr[i].Id, imgArr[i].URL)
	}
}

func downloadImage(fname string, url string) {
	imgpath := fmt.Sprintf("%s/image_original", dpath)

	err := os.MkdirAll(imgpath, os.ModePerm)

	if err != nil {
		fmt.Println("[Err] MkdirAll: ", err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.jpg", imgpath, fname))

	if err != nil {
		fmt.Println("[Err] Create File: ", err)
		return
	}

	defer file.Close()

	response, err := http.Get(url)

	if err != nil {
		fmt.Println("[Err] Get Response: ", err)
		return
	}

	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		fmt.Println("[Err] Copy File: ", err)
		return
	}

	fmt.Printf("[Info] Downloaded: %s\n", url)
}
