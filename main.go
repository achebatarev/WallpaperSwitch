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
	"time"
)

//TODO: create Tests
//TODO: manage image storage (cache, used before, limit size, etc.)
//TODO: manage how images are picked
//TODO: create a cli app

const root = "/tmp/wallpaper/"

func DownloadFile(url, name string) error {
	path := fmt.Sprint(root, name)
	//TODO: if folder is missing we ll be dead, fix it
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func DonwloadWallpaper() (string, error) {
	values := make(map[string]interface{})
	resp, err := http.Get("https://wallhaven.cc/api/v1/search")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &values)

	if err != nil {
		return "", err
	}


	data := values["data"].([]interface{})
	wallpaper := data[rand.Intn(len(data))].(map[string]interface{})
	link := wallpaper["path"].(string)
	id := wallpaper["id"].(string)

	//TODO: make sure that extension fits the file extension
	name := fmt.Sprintf("%s.png", id)

	err = DownloadFile(link, name)

	if err != nil {
		return "", err
	}

	return name, nil
}

func SwitchWallpaper(filename string) error {
	path := fmt.Sprint(root, filename)
	arg0 := "--bg-scale"
	cmd := exec.Command("feh", arg0, path)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())
	name, err := DonwloadWallpaper()

	if err != nil {
		log.Fatal(err)
	}

	err = SwitchWallpaper(name)

	if err != nil {
		log.Fatal(err)
	}

}
