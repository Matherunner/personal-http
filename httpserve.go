package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"html"
	"mime"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type EntryInfo struct {
	FileName string
	FileSize int64
}

var noDirList *bool
var dirEntryTempl *template.Template
var indexFile *string
var typesToGzip = []string{
	"text/plain", "text/html", "text/css", "application/javascript",
	"application/x-javascript", "text/javascript",
}

const dirListHeadFmt = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
h1 {
font-weight: normal;
font-size: large;
color: #777;
}
table {
border-collapse: collapse;
}
td, th {
border: 1px solid #ccc;
padding: 4px;
}
thead {
background: #f4f4f4;
color: #444;
text-shadow: 1px 1px #ddd;
}
.col-right {
text-align: right;
}
.col-left {
text-align: left;
padding-right: 3em;
}
tt {
white-space: pre;
}
a {
text-decoration: none;
}
a:hover {
text-decoration: underline;
}
</style>
<title>Directory listing of %[1]s</title>
</head>
<body>
<h1>%[1]s</h1>
<table>
<thead><tr><th class="col-left">name</th><th class="col-right">size</th></tr></thead>
<tbody>`

func useGzip(r *http.Request, mimeType string) bool {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}
	for _, v := range typesToGzip {
		if strings.Contains(mimeType, v) {
			return true
		}
	}
	return false
}

func outputDirList(w io.Writer, entryInfos []os.FileInfo, urlPath string) {
	fmt.Fprintf(w, dirListHeadFmt, html.EscapeString(urlPath))
	fmt.Fprint(w, `<tr><td class="col-left" colspan="2"><a href="../">[<i>parent dir</i>]</a></td></tr>`)
	for _, v := range entryInfos {
		entryInfo := EntryInfo{v.Name(), v.Size()}
		if v.IsDir() {
			entryInfo.FileName += "/"
		}
		dirEntryTempl.Execute(w, entryInfo)
	}
	fmt.Fprint(w, `</tbody></table></body></html>`)
}

func redirectSlash(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path[len(r.URL.Path) - 1] != '/' {
		http.Redirect(w, r, r.URL.Path + "/", http.StatusMovedPermanently)
		return true
	}
	return false
}

func serveDirList(w http.ResponseWriter, r *http.Request, fpath string) {
	if *indexFile != "" {
		indexPath := filepath.Join(fpath, *indexFile)
		if finfo, err := os.Stat(indexPath); err == nil {
			if !redirectSlash(w, r) {
				serveFile(w, r, indexPath, finfo)
			}
			return
		}
	}

	if *noDirList {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if redirectSlash(w, r) {
		return
	}

	entryInfos, err := ioutil.ReadDir(fpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read dir (%s)", err.Error()),
			http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		gzwriter := gzip.NewWriter(w)
		defer gzwriter.Close()
		outputDirList(gzwriter, entryInfos, r.URL.Path)
	} else {
		outputDirList(w, entryInfos, r.URL.Path)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, fpath string, finfo os.FileInfo) {
	file, err := os.Open(fpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot open file (%s)", err.Error()),
			http.StatusInternalServerError)
		return
	}
	defer file.Close()

	extType := mime.TypeByExtension(filepath.Ext(fpath))
	if useGzip(r, extType) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", extType)
		gzwriter := gzip.NewWriter(w)
		io.Copy(gzwriter, file)
		gzwriter.Close()
		return
	}

	// HACK: golang's ServeContent couldn't handle If-Range correctly if it
	// has Last-Modified times instead of ETag.
	r.Header.Del("If-Range")
	http.ServeContent(w, r, fpath, finfo.ModTime(), file)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s | %s %s %s\n", time.Now(),
		strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.Path)

	fpath := filepath.Join(".", r.URL.Path)
	finfo, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("cannot stat file (%s)", err.Error()),
			http.StatusInternalServerError)
		return
	}

	if finfo.IsDir() {
		serveDirList(w, r, fpath)
	} else {
		serveFile(w, r, fpath, finfo)
	}
}

func main() {
	dirEntryTempl, _ = template.New("Entry").Parse(`<tr><td class="col-left"><a href="{{.FileName}}"><tt>{{.FileName}}</tt></a></td><td class="col-right">{{.FileSize}}</td></tr>`)

	port := flag.Uint("p", 8080, "port number to listen")
	initPath := flag.String("path", ".", "working directory")
	noDirList = flag.Bool("nd", false, "disable directory listing")
	indexFile = flag.String("index", "", "default file to serve from a directory")
	flag.Parse()

	*indexFile = strings.TrimSpace(*indexFile)
	os.Chdir(*initPath)
	http.HandleFunc("/", mainHandler)
	portStr := strconv.FormatUint(uint64(*port), 10)
	fmt.Printf("Listening on %s...\n", portStr)
	err := http.ListenAndServe(":" + portStr, nil)
	if err != nil {
		fmt.Printf("Failed to listen (%s)\n", err.Error())
	}
}
