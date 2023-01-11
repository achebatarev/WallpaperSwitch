package download

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/achebatarev/WallpaperSwitch/config"
)

type Wallpaper struct {
	Id        string
	Extension string
	Link      string
	Name      string
	Preview   string
}

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
	path := fmt.Sprint(config.Root, wallpaper.Name)

	if err := os.MkdirAll(config.Root, os.ModePerm); err != nil {
		return fmt.Errorf("DownloadFile: Could not create folder: %w", err)
	}

	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("DownloadFile: Could not create file: %w", err)
	}

	resp, err := http.Get(wallpaper.Link)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("DownloadFile: Could not get url %q: %w", wallpaper.Link, err)
	}

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
