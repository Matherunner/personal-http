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

#maindiv {
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
</script>
</head>
<body>
<div id="maindiv">
<form action="/upload" enctype="multipart/form-data" method="post">
<input type="file" name="file" onchange="handlefile(this.form);" style="display:none">
<input class="mainbtn" type="button" name="select" value="Select a file" onclick="selfile(this.form);">
<input class="mainbtn" type="submit" name="upload" value="Send" disabled="disabled">
</form>
</div>
</body>
</html>
