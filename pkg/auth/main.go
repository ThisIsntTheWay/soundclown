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
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
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

	driver.Get("https://soundcloud.com/signin")

	// Feb 2023
	// email: sign_in_up_email     <ID, input>
	// pass:  enter_password_field <ID, input>
	fmt.Println("Attempting sign in...")
	elem, err := driver.FindElement(selenium.ByID, "sign_in_up_email")
	if err != nil {
		panic(err)
	}

	// TODO, put into text
	elem.

	return err
}
