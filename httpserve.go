package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EntryInfo struct {
	Name string
	Size int64
}

var storePath *string
var rootPath *string
var tmplListHead *template.Template
var tmplListEntry *template.Template
var tmplListFoot *template.Template

func logConnection(r *http.Request) {
	log.Printf("[%s] %s %s\n", r.RemoteAddr, r.Method, r.URL)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	logConnection(r)
	w.WriteHeader(http.StatusNotFound)
}

func serveDirList(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.URL.Path[len(r.URL.Path) - 1] != '/' {
		http.Redirect(w, r, r.URL.Path + "/", http.StatusMovedPermanently)
		return
	}

	fileInfos, err := ioutil.ReadDir(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplListHead.Execute(w, "Listing of " + strings.TrimPrefix(r.URL.Path, "/files"))
	for _, v := range fileInfos {
		entryInfo := EntryInfo{v.Name(), v.Size()}
		if v.IsDir() {
			entryInfo.Name += "/"
		}
		tmplListEntry.Execute(w, entryInfo)
	}
	tmplListFoot.Execute(w, nil)
}

func serveFile(w http.ResponseWriter, r *http.Request, filePath string, fileInfo os.FileInfo) {
	file, err := os.Open(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.ServeContent(w, r, filePath, fileInfo.ModTime(), file)
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	logConnection(r)
	filePath := strings.TrimPrefix(r.URL.Path, "/files")
	filePath = filepath.Join(*rootPath, filePath)
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		serveDirList(w, r, filePath)
	} else {
		serveFile(w, r, filePath, fileInfo)
	}
}

func serveUpload(w http.ResponseWriter, r *http.Request) {
	logConnection(r)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("upload.gtpl")
		if err != nil {
			log.Fatalf("Failed to parse upload.gtpl (%s)\n", err.Error())
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("[%s] Failed to write upload form to client (%s)\n:",
				r.RemoteAddr, err.Error())
			return
		}
	} else if r.Method == "POST" {
		w.WriteHeader(http.StatusNoContent)
		mpReader, err := r.MultipartReader()
		if err != nil {
			log.Printf("[%s] Failed to read multipart form (%s)\n",
				r.RemoteAddr, err.Error())
			return
		}

		log.Printf("[%s] Got multipart form of size %d\n",
			r.RemoteAddr, r.ContentLength)

		var part *multipart.Part
		for {
			part, err = mpReader.NextPart()
			if err == io.EOF {
				return
			} else if err != nil {
				continue
			}

			if part.FileName() != "" {
				log.Printf("[%s] Downloading \"%s\"\n", r.RemoteAddr, part.FileName())
				break
			}

			part.Close()
		}
		defer part.Close()

		fileName := filepath.Join(*storePath, filepath.Base(part.FileName()))
		outFile, err := os.OpenFile(fileName, os.O_WRONLY | os.O_CREATE | os.O_EXCL,
			0600)
		if err != nil {
			log.Printf("[%s] Failed to create file (%s)\n", r.RemoteAddr, err.Error())
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, part)
		if err != nil {
			log.Printf("[%s] Error occured when transferring \"%s\" (%s)\n",
				r.RemoteAddr, part.FileName(), err.Error())
			return
		}

		log.Printf("[%s] Transfer of \"%s\" completed\n", r.RemoteAddr, part.FileName())
	}
}

func main() {
	portNum := flag.Uint("p", 8080, "port number to open")
	storePath = flag.String("store", ".", "path to store downloaded files")
	rootPath = flag.String("root", ".", "root directory for serving files")
	flag.Parse()

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/files/", serveFiles)
	http.HandleFunc("/upload", serveUpload)

	var err error
	tmplListHead, err = template.ParseFiles("filelist-head.gtpl")
	if err != nil {
		log.Fatalf("Failed to parse filelist-head.gtpl (%s)", err.Error())
	}

	tmplListEntry, err = template.ParseFiles("filelist-entry.gtpl")
	if err != nil {
		log.Fatalf("Failed to parse filelist-entry.gtpl (%s)", err.Error())
	}

	tmplListFoot, err = template.ParseFiles("filelist-foot.gtpl")
	if err != nil {
		log.Fatalf("Failed to parse filelist-foot.gtpl (%s)", err.Error())
	}

	log.Printf("Listening to port %d\n", *portNum)
	err = http.ListenAndServe(":" + strconv.FormatUint(uint64(*portNum), 10), nil)
	if err != nil {
		log.Fatalf("Failed to listen to port (%s)\n", err.Error())
	}
}
