package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raff/godet"
)

func main() {

	// currentURL := flag.Bool("current", false, "print current URL open")
	// openURL := flag.String("open", "", "URL to open in browser")

	// flag.Parse()

	// user requested the current open URL
	// if *currentURL {
	// 	tabs, _ := remoteConn.TabList("")
	// 	for _, value := range tabs {
	// 		if value.Type != "background_page" {
	// 			fmt.Printf("%s,%s\n", value.Title, value.URL)
	// 		}
	// 	}
	// }

	// user wants to open a new URL
	// if *openURL != "" {
	// 	_, _ = remoteConn.Navigate(*openURL)
	// }

	router := mux.NewRouter()
	router.HandleFunc("/", printCurrentURL).Methods("GET")
	router.HandleFunc("/current", printCurrentURL).Methods("GET")
	router.HandleFunc("/open", openURLInBrowser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func printCurrentURL(w http.ResponseWriter, r *http.Request) {

	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
		return
	}

	// disconnect when done
	defer remote.Close()
	tabs, _ := remote.TabList("")
	for _, value := range tabs {
		if value.Type != "background_page" {
			fmt.Fprintln(w, value.Title, value.URL)
		}
	}
}

func openURLInBrowser(w http.ResponseWriter, r *http.Request) {

	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		fmt.Fprintln(w, "cannot connect to Chrome instance:", err)
		return
	}

	// disconnect when done
	defer remote.Close()

	vars := mux.Vars(r)
	openURL := vars["url"]
	fmt.Println(openURL)
	_, _ = remote.Navigate(openURL)
}
