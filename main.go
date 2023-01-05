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

//TODO: create Tests
//TODO: manage image storage (cache, used before, limit size, etc.)
//TODO: manage how images are picked
//TODO: create a cli app
// TODO: parse response into a list of structs
//Idea: can use a package with cross-platform support for changing wallpapers
// then we just focus on creating a cli app to interact with wallpaper selection
// defintely soething worth implementing down the road

type Wallpaper struct{
    id string
    extension string
    link string
    name string
}

const root = "/tmp/wallpaper/"


func Parse(object []byte) ([]Wallpaper, error){
    wallpapers := []Wallpaper{}
    values := make(map[string]interface{})

    if err := json.Unmarshal(object, &values); err != nil{
        return nil, fmt.Errorf("Parse: %w", err)
    }

    data, ok := values["data"].([]interface{})

    if !ok{
        return nil, fmt.Errorf("Parse: data is not []interface{}")
    }
    
    for _, e := range data{
        wallpaper := e.(map[string]interface{})

        id := wallpaper["id"].(string)
        link := wallpaper["path"].(string)
        filetype := wallpaper["file_type"].(string)
        extension := strings.Split(filetype, "/")[1]
        name := fmt.Sprintf("%s.%s", id, extension)
        new_wallpaper := Wallpaper{id, extension, link, name}

        wallpapers = append(wallpapers, new_wallpaper)
    }

    return wallpapers, nil
    
}

func DownloadFile(wallpaper Wallpaper) error {
	path := fmt.Sprint(root, wallpaper.name)
	//TODO: if folder is missing we ll be dead, fix it
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

func DonwloadWallpaper() (string, error) {
	resp, err := http.Get("https://wallhaven.cc/api/v1/search?categories=010")

	if err != nil {
        return "", fmt.Errorf("DownloadWallpaper: %w", err) 
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
        return "", fmt.Errorf("DownloadWallpaper: %w", err) 
	}

    wallpapers, _ := Parse(body)
	wallpaper := wallpapers[rand.Intn(len(wallpapers))]

	if err := DownloadFile(wallpaper); err != nil {
		return "", fmt.Errorf("DownloadWallpaper: %w", err)
	}

	return wallpaper.name, nil
}

func SwitchWallpaper(filename string) error {
	path := fmt.Sprint(root, filename)
	arg0 := "--bg-scale"
	cmd := exec.Command("feh", arg0, path)

	if _, err := cmd.Output(); err != nil{
        return fmt.Errorf("SwitchWallpaper %w", err)
    }

	return nil
}

func main() {
	rand.Seed(time.Now().Unix())
	name, err := DonwloadWallpaper()

	if err != nil {
		log.Fatal(err)
	}

    if err := SwitchWallpaper(name); err != nil{
		log.Fatal(err)
    }


}
