<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Send a file</title>
<style>
.mainbtn {
  display: block;
  font-size: 12pt;
  margin: 10px auto 0 auto;
  outline: none;
  padding: 10px;
  width: 100%;
  background-color: #dedede;
  border: 0;
}

#progress {
  display: none;
  width: 100%;
  margin: 10px auto 0 auto;
}

.mainbtn::-moz-focus-inner {
  border: 0;
}

.mainbtn:hover {
  background-color: #ccc;
}

.mainbtn:active {
  background-color: #afafaf;
}

.mainbtn:disabled {
  background-color: #f3f3f3;
  color: #aaa;
}

#mainform {
  position: relative;
  margin: 0 auto 0 auto;
  max-width: 1000px;
  font-family: sans-serif;
}
</style>
<script>
function selfile(f) {
  f.file.click();
}

function handlefile(f) {
  f.select.value = f.file.files[0].name;
  f.upload.disabled = false;
}

function uploadfile(f) {
  var xhr = new XMLHttpRequest();
  var prog = document.getElementById("progress");
  prog.style.display = "block";
  f.upload.disabled = true;

  xhr.upload.addEventListener("progress", function(e) {
    prog.value = e.loaded / e.total;
  }, false);

  xhr.addEventListener("readystatechange", function(e) {
    if (xhr.readyState === 4) {
      prog.style.display = "none";
      prog.value = 0;
      f.upload.disabled = false;
    }
  }, false);

  var fd = new FormData();
  xhr.open("POST", "/upload", true);
  fd.append("file", f.file.files[0]);
  xhr.send(fd);
}
</script>
</head>
<body>
<form id="mainform">
<input type="file" name="file" onchange="handlefile(this.form);" style="display:none">
<input class="mainbtn" type="button" name="select" value="Select a file" onclick="selfile(this.form);">
<input class="mainbtn" type="button" name="upload" value="Send" onclick="uploadfile(this.form);" disabled="disabled">
<progress id="progress" value="0"></progress>
</form>
</body>
</html>
