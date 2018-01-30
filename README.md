# Hobbit

## DESCRIPTION

Extremely simple console tool. Takes Long URL as a first argument.
Returns Bitly short URL to console. That's It.

## PREREQUISITES

* OS: Linux/macOS. Most likely will work on Windows as well but not tested.
* Bitly account. Get it here <https://bitly.com>

## BUILD

* Get and install Go <https://golang.org/dl/>
* Run `go build hobbit.go`

## INSTALL

* Build binary or download It from [releases](https://github.com/shuienko/hobbit/releases)
* Rename file to `hobbit`
* Make It executable: `chmod +x hobbit`
* Place `hobbit` binary under your `PATH`

## USAGE

* On first run app will ask for Bitly credentials in order to obtain API Token.
* `hobbit http://example.com`
* Please keep in mind that `http`/`https` part is obligatory
