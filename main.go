package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Wallpaper struct {
	id        string
	extension string
	link      string
	name      string
	preview   string
}

const root = "/tmp/wallpaper/"

func Parse(object []byte) ([]Wallpaper, error) {
	wallpapers := []Wallpaper{}
	values := make(map[string]interface{})

	if err := json.Unmarshal(object, &values); err != nil {
		return nil, fmt.Errorf("Parse: %w", err)
	}

	data, ok := values["data"].([]interface{})

	if !ok {
		return nil, fmt.Errorf("Parse: data is not []interface{}")
	}

	for _, e := range data {
		wallpaper := e.(map[string]interface{})

		id := wallpaper["id"].(string)
		link := wallpaper["path"].(string)
		filetype := wallpaper["file_type"].(string)
		extension := strings.Split(filetype, "/")[1]
		name := fmt.Sprintf("%s.%s", id, extension)
		thumb := wallpaper["thumbs"].(map[string]interface{})["large"].(string)

		new_wallpaper := Wallpaper{id, extension, link, name, thumb}

		wallpapers = append(wallpapers, new_wallpaper)
	}

	return wallpapers, nil

}

func DownloadFile(wallpaper *Wallpaper) error {
	path := fmt.Sprint(root, wallpaper.name)

	if err := os.MkdirAll(root, os.ModePerm); err != nil {
		return fmt.Errorf("DownloadFile: Could not create folder: %w", err)
	}

	file, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("DownloadFile: Could not create file: %w", err)
	}

	defer file.Close()
	resp, err := http.Get(wallpaper.link)

	if err != nil {
		return fmt.Errorf("DownloadFile: Could not get url %q: %w", wallpaper.link, err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("DownloadFile: Writing to file failed: %w", err)
	}

	return nil
}

func DonwloadWallpaper() (*Wallpaper, error) {
	resp, err := http.Get("https://wallhaven.cc/api/v1/search?categories=010")

	if err != nil {
		return nil, fmt.Errorf("DownloadWallpaper: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("DownloadWallpaper: %w", err)
	}

	wallpapers, _ := Parse(body)
	wallpaper := &wallpapers[rand.Intn(len(wallpapers))]

	if err := DownloadFile(wallpaper); err != nil {
		return nil, fmt.Errorf("DownloadWallpaper: %w", err)
	}

	return wallpaper, nil
}

func SwitchWallpaper(wallpaper *Wallpaper) error {
	path := fmt.Sprint(root, wallpaper.name)
	arg0 := "--bg-scale"
	cmd := exec.Command("feh", arg0, path)

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("SwitchWallpaper %w", err)
	}

	return nil
}

func DisplayWallpaper(wallpaper *Wallpaper) (*exec.Cmd, error) {
	arg0 := "-x"
	cmd := exec.Command("feh", arg0, wallpaper.preview)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("DisplayWallpaper %w", err)
	}

	return cmd, nil
}

func HandleInput() bool {
	var input string

	for true {
		fmt.Print("Do you want to set this image? (y/n) ")
		fmt.Scanln(&input)
		switch input {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Please select between y and n")
		}
	}

	return false
}

func main() {
	rand.Seed(time.Now().Unix())
	var wallpaper *Wallpaper

	done := false
	for !done {
		var err error

		wallpaper, err = DonwloadWallpaper()

		if err != nil {
			log.Fatal(err)
		}

		cmd, err := DisplayWallpaper(wallpaper)

		if err != nil {
			log.Fatal(err)
		}

		done = HandleInput()

		if err := cmd.Process.Kill(); err != nil {
			log.Fatal(err)
		}

	}

	if err := SwitchWallpaper(wallpaper); err != nil {
		log.Fatal(err)
	}

}
