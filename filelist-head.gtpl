<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>{{.}}</title>
<style>
body {
  margin: 15px;
}

a {
  text-decoration: none;
  color: #133170;
}

a:hover {
  text-decoration: underline;
}

table {
  border-collapse: collapse;
}

thead {
  font-weight: bold;
}

tbody {
  margin-top: 100px;
  padding-top: 100px;
}

tr {
  border-bottom: 1px solid #ddd;
}

td.nam {
  padding-right: 30px;
  padding-left: 7px;
}

td.siz {
  text-align: right;
  padding-right: 7px;
}

tr.spacer {
  height: 10px;
}

tr.head-row {
  border: none;
}
</style>
</head>
<body>
<table>
<thead>
<tr class="head-row"><td class="nam">Name</td><td class="siz">Size</td></tr>
</thead>
<tbody>
<tr class="spacer"></tr>
