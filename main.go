package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

//go:embed upload.html
var uploadHtml string
var downloadDir = ""
var uploadDir = ""
var listenPort = 80

func init() {
	flag.StringVar(&downloadDir, "d", "", "download from dir")
	flag.StringVar(&uploadDir, "u", "", "upload to dir")
	flag.IntVar(&listenPort, "p", 80, "listen port")
	flag.Parse()
}

func main() {
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(downloadDir))))
	http.HandleFunc("/", UploadHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(uploadHtml))
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		f, _, err := r.FormFile("file")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		fileName := ""
		files := r.MultipartForm.File
		for k, v := range files {
			for _, vv := range v {
				fmt.Println(k, ":", vv.Filename)
				fileName = vv.Filename
				break
			}
		}

		saveFile, err := os.Create(path.Join(uploadDir, fileName))
		if err != nil {
			log.Fatalf("Create file: %v", err)
		}
		defer saveFile.Close()
		io.Copy(saveFile, f)

		w.Write([]byte("File Uploaded Successfully"))
		return
	}
}
