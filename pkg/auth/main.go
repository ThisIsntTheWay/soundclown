package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// PUBLIC

func GetAuthToken() string {
	return "abc"
}

func main() {
	// Debug stuff to get current filepath
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	fmt.Println("Eecutable working dir: " + exPath)

	startSelenium()
}

// PRIVATE

func startSelenium() error {
	// Run Chrome browser
	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		// "--headless",  // comment out this line to see the browser
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}

	driver.Get("https://www.google.com")

	return err
}
