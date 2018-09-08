package main

import (
	"flag"
	"fmt"

	"github.com/raff/godet"
)

func main() {

	currentURL := flag.Bool("current", false, "print current URL open")
	openURL := flag.String("open", "", "URL to open in browser")

	flag.Parse()

	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		fmt.Println("cannot connect to Chrome instance:", err)
		return
	}

	// disconnect when done
	defer remote.Close()

	// user requested the current open URL
	if *currentURL {
		tabs, _ := remote.TabList("")
		for _, value := range tabs {
			if value.Type != "background_page" {
				fmt.Printf("%s,%s\n", value.Title, value.URL)
			}
		}
	}

	// user wants to open a new URL
	if *openURL != "" {
		_, _ = remote.Navigate(*openURL)
	}

}
