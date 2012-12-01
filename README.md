munin
============

Client library for munin in go

LICENSE
-------
BSD

documentation
-------------
[package documentation at go.pkgdoc.org](http://go.pkgdoc.org/github.com/nightlyone/munin)

compile and install
-------------------
Get the library
	go get github.com/nightlyone/munin

Get the mfgo client
	go get github.com/nightlyone/munin/cmd/mfgo

Play around with the client
	$GOROOT/bin/mfgo -h

contributing
============

Contributions are welcome. Please open an issue or send me a pull request for a dedicated branch.
Make sure the git commit hooks show it works.

git commit hooks
-----------------------
enable commit hooks via

        cd .git ; rm -rf hooks; ln -s ../git-hooks hooks ; cd ..

