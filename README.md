# Remote Control of Chromium

## Depends

This program depends upon, https://github.com/raff/godet.  This library does the hard work of talking to the remote debug port of Chrome.

## Requirements

You must be using Chromium (maybe Chrome) and you must start it with remote debugging turned on.  The included shell script will start the browser with remote debugging and sets a temporary user profile.

You must fill in the necessary items in the config.toml file.

### Building

Use the appropriate build_ARCH.sh script

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

## Config TOML Format

### Listen Section

This section is used to configure the port you want the web interface for this tool to listen.  It is also where you configure SSL support.

```shell
[listen]
ssl = false
cert = "wildcard_certificate"
key = "wildcard_key"
port = 8081
```

### Chrome Section

* `host` should be localhost but it is the network address the Chromium remote debugging is configured for.

* `port` is the default Chromium remote debugging port

```shell
[chrome]
host = "localhost"
port = 9222
```

## Example Starting Chromium

The following is a shell script that shows how Chromium could be started.  Notice `PREVURL`.  The program will write out the URL that was sent to the client so that it can be loaded by default.  The location of the file is $HOME/urlfile.txt.  This is on line 82 of main.go.

```shell
#!/bin/sh
CHROME_DATA_DIR=$(mktemp -d)
trap "rm -rf ${CHROME_DATA_DIR}" SIGINT SIGTERM EXIT

PREVURL=$(cat /home/pi/urlfile.txt)
DEFAULTURL="https://google.com"
URL=${PREVURL:-$DEFAULTURL}

/usr/bin/chromium-browser --remote-debugging-port=9222 --user-data-dir="${CHROME_DATA_DIR}" --disable-infobars --kiosk "${URL}"
```
