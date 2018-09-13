# Remote Control of Chromium

## Depends

This program depends upon, https://github.com/raff/godet.  This library does the hard work of talking to the remote debug port of Chrome.

## Requirements

You must be using Chromium (maybe Chrome) and you must start it with remote debugging turned on.  The included shell script will start the browser with remote debugging and sets a temporary user profile.

You must fill in the necessary items in the config.toml file.

### Server

* The listen port for the adminchrome web server
* If you want to enable SSL or not
* The path to the necessary SSL certs if you want SSL

### Chrome

* The remote debugging address.  This is best left to localhost
* The remote debugging port.  You can leave this as the default

## Why

The idea is that this program is the basis for remotely handling controlling a browser that would be used as a digital sign.

## Launching

```shell
$ ./remotechrome -help
Usage of ./remotechrome:
  -conf
        Config file for this listener and chrome port info
```

## Usage

### View the current page being viewed by Chromium

```shell
curl http://localhost:9222/
or
curl http://localhost:9222/current
```

### Open a URL

```shell
curl -X POST -d "https://www.google.com" http://localhost:9222/open
```


