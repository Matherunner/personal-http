# httpserve

Have you ever been frustrated because even in 2015

- you couldn't transfer huge files easily between two computers on a LAN

- you couldn't easily send huge files to someone

- you couldn't get someone to easily send huge files to you

then `httpserve` is a solution you might consider.  Note the keyword _easily_.
Sure, you could transfer relatively small files through email attachment, but
good luck sending a _10 GB file_ this way.  Don't even talk about cloud
services like Google Drive and Dropbox, I shouldn't have to upload and store a
large file somewhere only to have it downloaded and deleted later.  The same
goes for external USB storage.

Now you might be a geek and have a web server to host any file at your
disposal.  That's great for you, but I can assure you that most people will not
welcome the hassle of setting up a http daemon just to transfer some darn
files.  You might also be quick to remind me of solutions such as

    $ python -mhttp.server 8080
    $ twistd -n web -port=8080 -path=.

While they may work fine for you to send files to someone, they do not solve
the problem of getting someone (especially non-geeks) to send files to you.
You might then say you could write a file upload form in PHP and run it through
Apache.  But the point is nobody wants to install Apache, PHP or some other
heavy weight software and edit a bunch of config files just to do file
transfer.  And not everyone has as much free time and energy.

As for transferring files between two computers on a LAN, you might be eager to
"inform" me how I could spend time setting up and configure NFS or Samba on
both computers.  But why the hassle?  I just want to transfer some files and
call it a day!

## Usage

With `httpserve` you just need an open port on your side, then you can transfer
files to someone or have someone send you files, through HTTP.

Say you have two computers, A (192.168.0.1) and B (192.168.0.2), separated by a
wireless link.  You want to send a 10 GB file from A to B.  If `httpserve` is
installed on A, then you can run

    $ httpserve -p 9090 -root /home/me/docs

to serve the files in `/home/me/docs` via port 9090.  On B you will point the
browser to `http://192.168.0.1:9090/files`.  A directory listing will be
produced, with hyperlinks to each file.  Suppose `httpserve` is installed on B
instead.  You would run

    $ httpserve -p 9090 -store /home/me/docs

to send a minimalistic upload form to A via port 9090.  Then select the 10 GB
file on A and send it.  The file will be stored in `/home/me/docs`.

## How to build

You will need to have [golang](https://golang.org/) installed.  Then with

    $ go build httpserve.go

and you're good to go.
