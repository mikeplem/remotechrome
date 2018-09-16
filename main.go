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
	"time"

	"github.com/BurntSushi/toml"
	"github.com/raff/godet"
)

func printCurrentURL(w http.ResponseWriter, r *http.Request) {

	connString := fmt.Sprintf("%s:%d", Config.Chrome.Host, Config.Chrome.Port)

	remote, err := godet.Connect(connString, false)
	if err != nil {
		fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
		return
	}

	defer remote.Close()

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
			log.Fatal(err)
		}
		unescapeURL, err := url.PathUnescape(u.String())
		if err != nil {
			log.Println(err)
		}

		stripURL := strings.TrimPrefix(unescapeURL, "u=")

		remote, err := godet.Connect(connString, false)
		if err != nil {
			fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
			return
		}

		defer remote.Close()

		fmt.Printf("Requested to open %s\n", stripURL)
		_, _ = remote.Navigate(stripURL)
		fmt.Fprint(w, "POST done")

		urlToFile := []byte(stripURL)

		fileToWrite := fmt.Sprintf("%s/urlfile.txt", userHome)

		err = ioutil.WriteFile(fileToWrite, urlToFile, 0644)
		if err != nil {
			log.Println(err)
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

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

	if Config.Listen.SSL == true {
		fmt.Println("Listening on port " + listenPort + " with SSL")
		err := http.ListenAndServeTLS(listenPort, Config.Listen.Cert, Config.Listen.Key, nil)
		if err != nil {
			fmt.Print(time.Now())
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		fmt.Println("Listening on port " + listenPort + " without SSL")
		err := http.ListenAndServe(listenPort, nil)
		if err != nil {
			fmt.Print(time.Now())
			log.Fatal("ListenAndServe: ", err)
		}
	}

}
