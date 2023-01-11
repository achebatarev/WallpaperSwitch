package main

import (
	"log"
	"math/rand"
	"time"
)

// TODO: add an ability for users to change the filename of a file saved in favorites
// TODO: add Parsing of a config file
// TODO: Check for a config file using  environment variables
// TODO: set enviroment variable for a config file same as Home Directory

const root = "/tmp/wallpaper/"
const favorite = "/home/alex/Pictures/Wallpapers/"

func main() {
	rand.Seed(time.Now().Unix())
	if err := TUI(); err != nil {
		log.Fatal(err)
	}

}
