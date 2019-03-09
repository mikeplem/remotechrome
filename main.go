package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/raff/godet"
)

// User Home directory
var userHome = os.Getenv("HOME")

// ConfigFile holds the user supplied configuration file - it is placed here since it is a global
var ConfigFile *string

// Config is the structure of the TOML config structure
var Config tomlConfig

type tomlConfig struct {
	Listen listenconfig `toml:"listen"`
	Chrome chromeconfig `toml:"chrome"`
}

type listenconfig struct {
	SSL  bool
	Cert string
	Key  string
	Port int
}

type chromeconfig struct {
	Host string
	Port int
}

func viewScreenshot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "screenshot.png")
}

func saveScreenshot(w http.ResponseWriter, r *http.Request) {

	connString := fmt.Sprintf("%s:%d", Config.Chrome.Host, Config.Chrome.Port)

	remote, err := godet.Connect(connString, false)
	if err != nil {
		fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
		return
	}

	defer remote.Close()

	log.Print("Taking screenshot")
	err = remote.SaveScreenshot("screenshot.png", 0644, 0, true)
	if err != nil {
		log.Print(err)
		return
	}
}

func reloadBrowser(w http.ResponseWriter, r *http.Request) {

	connString := fmt.Sprintf("%s:%d", Config.Chrome.Host, Config.Chrome.Port)

	remote, err := godet.Connect(connString, false)
	if err != nil {
		fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
		return
	}

	defer remote.Close()

	log.Print("Reloading browser page")
	err = remote.Reload()
	if err != nil {
		log.Print(err)
		return
	}
}

func printCurrentURL(w http.ResponseWriter, r *http.Request) {

	connString := fmt.Sprintf("%s:%d", Config.Chrome.Host, Config.Chrome.Port)

	remote, err := godet.Connect(connString, false)
	if err != nil {
		fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
		return
	}

	defer remote.Close()

	// when the browser starts there is a hidden tab called background_page.
	// do not show that in the list of open tabs
	tabs, _ := remote.TabList("")
	for _, value := range tabs {
		if value.Type != "background_page" {
			fmt.Fprintln(w, value.Title, value.URL)
		}
	}
}

func openURLInBrowser(w http.ResponseWriter, r *http.Request) {

	var results []string
	var openURL string

	connString := fmt.Sprintf("%s:%d", Config.Chrome.Host, Config.Chrome.Port)

	// 	helpful code https://gist.github.com/alyssaq/75d6678d00572d103106

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}

		results = append(results, string(body))
		openURL = strings.Join(results, " ")

		u, err := url.Parse(openURL)
		if err != nil {
			errorHandler(w, r, http.StatusBadRequest)
			return
		}

		unescapeURL, err := url.PathUnescape(u.String())
		if err != nil {
			log.Print(err)
		}

		stripURL := strings.TrimPrefix(unescapeURL, "u=")

		remote, err := godet.Connect(connString, false)
		if err != nil {
			fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
			return
		}

		defer remote.Close()

		log.Printf("Requested to open %s\n", stripURL)
		_, _ = remote.Navigate(stripURL)
		fmt.Fprint(w, "POST done")

		urlToFile := []byte(stripURL)

		fileToWrite := fmt.Sprintf("%s/urlfile.txt", userHome)

		err = ioutil.WriteFile(fileToWrite, urlToFile, 0644)
		if err != nil {
			log.Print(err)
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	switch status {
	case http.StatusGone:
		log.Print("The command does not exist.")
		return
	case http.StatusBadRequest:
		log.Print("The URL did not open for some reason.")
		return
	}

}

func init() {

	ConfigFile = flag.String("conf", "", "Config file for this listener and chrome port info")

	flag.Parse()

	if _, err := toml.DecodeFile(*ConfigFile, &Config); err != nil {
		log.Fatal(err)
	}

}

func main() {

	listenPort := fmt.Sprintf(":%d", Config.Listen.Port)

	http.HandleFunc("/", printCurrentURL)
	http.HandleFunc("/current", printCurrentURL)
	http.HandleFunc("/open", openURLInBrowser)
	http.HandleFunc("/reload", reloadBrowser)
	http.HandleFunc("/screenshot", saveScreenshot)
	http.HandleFunc("/view", viewScreenshot)

	if Config.Listen.SSL == true {
		log.Print("Listening on port " + listenPort + " with SSL")
		err := http.ListenAndServeTLS(listenPort, Config.Listen.Cert, Config.Listen.Key, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		log.Print("Listening on port " + listenPort + " without SSL")
		err := http.ListenAndServe(listenPort, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}

}
