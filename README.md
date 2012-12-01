munin
============

Client library for munin in go

[![Build Status][1]][2]

[1]: https://secure.travis-ci.org/nightlyone/munin.png
[2]: http://www.travis-ci.org/nightlyone/munin


LICENSE
-------
BSD

documentation
-------------
[package documentation at go.pkgdoc.org](http://go.pkgdoc.org/github.com/nightlyone/munin)

compile and install
-------------------
Install [Go 1][3], either [from source][4] or [with a prepackaged binary][5].
[3]: http://golang.org
[4]: http://golang.org/doc/install/source
[5]: http://golang.org/doc/install

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

