package main

import (
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const mainHtml = `<!doctype html><html>
<head>
<meta charset="utf-8">
<style>
html {
min-height: 100%;
text-align: center;
}
input[type=file] {
display: none;
}
.mybtn {
background: #e8e8e8;
border-radius: 4px;
border: 1px solid #b1b1b1;
box-shadow: 2px 2px 5px #dcdcdc;
color: #444;
display: block;
font-size: small;
height: 50px;
margin: 15px auto;
max-width: 1000px;
outline: none;
width: 100%;
}
.mybtn:hover:enabled {
border: 1px solid #888;
color: #111;
}
.mybtn:focus {
border: 1px solid #4064bf;
}
.mybtn:disabled {
opacity: 0.3;
}
.uploadbtn {
font-weight: bold;
}
</style>
<script>
function setup() {
	var fbtn = document.upform.filebtn, felem = document.upform.filename, subbtn = document.upform.submit;
	fbtn.addEventListener("click", function(e) { felem.click() }, false);
	felem.addEventListener("change", function(e) {
		if (felem.files.length) {
			fbtn.value = felem.files[0].name;
			subbtn.disabled = false;
		} else {
			fbtn.value = "select file";
			subbtn.disabled = true;
		}
	}, false);
}
window.addEventListener("load", setup, true);
</script>
<title>File upload</title>
</head>
<body>
<form method="post" action="/" enctype="multipart/form-data" name="upform">
<input type="file" name="filename">
<input type="button" class="mybtn" name="filebtn" value="select file">
<input type="submit" class="mybtn uploadbtn" name="submit" value="upload" disabled="true">
</form>
</body>`

func mainHandler(w http.ResponseWriter, r *http.Request) {
	remoteIPAddr := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Printf("%s | %s %s | %s\n", time.Now(), remoteIPAddr, r.Method,
		r.Header.Get("Content-Length"))

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if r.Method == "GET" {
		fmt.Fprint(w, mainHtml)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	if r.Method != "POST" {
		return
	}

	mpreader, err := r.MultipartReader()
	if err != nil {
		fmt.Printf("  Failed to read from %s (%s)\n", remoteIPAddr,
			err.Error())
		return
	}

	var filepart *multipart.Part
	for {
		filepart, err = mpreader.NextPart()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("  %s did not upload any file.\n", remoteIPAddr)
				return
			}
			continue
		}
		if filepart.FileName() != "" {
			fmt.Printf("  Got \"%s\" from %s\n", filepart.FileName(),
				remoteIPAddr)
			break
		}
		filepart.Close()
	}
	defer filepart.Close()

	// filepath.Base prevents saving the file in potentially any directory if
	// the supplied file name uses relative paths.
	outfile, err := os.OpenFile(filepath.Base(filepart.FileName()),
		os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0600)
	if err != nil {
		fmt.Printf("  Failed to create output file for %s's \"%s\" (%s)\n",
			remoteIPAddr, filepart.FileName(), err.Error())
		return
	}
	defer outfile.Close()

	_, err = io.Copy(outfile, filepart)
	if err != nil {
		fmt.Printf("  Error occurred while saving %s's \"%s\" (%s)\n",
			remoteIPAddr, filepart.FileName(), err.Error())
		return
	}
}

func main() {
	portNum := flag.Uint("p", 8080, "port number to listen")
	workingPath := flag.String("path", ".", "directory to save files in")
	flag.Parse()

	os.Chdir(*workingPath)
	http.HandleFunc("/", mainHandler)
	portNumStr := strconv.FormatUint(uint64(*portNum), 10)
	fmt.Printf("Listening on %s...\n", portNumStr)
	err := http.ListenAndServe(":" + portNumStr, nil)
	if err != nil {
		fmt.Printf("Failed to listen (%s)\n", err.Error())
	}
}
