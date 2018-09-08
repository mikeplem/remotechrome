# Remote Control of Chromium

## Depends

This program depends upon, https://github.com/raff/godet.  This library does the hard work of talking to the remote debug port of Chrome.

## Requirements

You must be using Chromium (maybe Chrome) and you must start it with remote debugging turned on.  The included shell script will start the browser with remote debugging and sets a temporary user profile.

## Why

The idea is that this program is the basis for remotely handling controlling a browser that would be used as a digital sign.

## Usage
I was looking for a simple way to remotely open a URL in Chromium.  The Go code takes two arguments.

```shell
$ ./remotechrome -help
Usage of ./remotechrome:
  -current
        print current URL open
  -open string
        URL to open in browser
```

## TODO

* Add REST endpoints
* Add LDAP authentication
