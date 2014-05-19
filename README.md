# Personal HTTP servers

`httpserve` is a simple web server which can output directory listings and serve static files.  It gzips the output for certain MIME types.  Use this as a quick way to transfer files, instead of the basic web servers provided by Python, Ruby or Twisted.

`httpup` does file uploads only.  It always displays a minimalistic upload interface to the end user.  This program can be used if you want a non-technical person to send files to you.

Obviously, both programs require an open port to work properly.  Be sure to configure your firewall and NAT router.

To build these programs, just do

    go build httpserve.go
    go build httpup.go
