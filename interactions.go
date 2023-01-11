package main

import (
	"fmt"

	"github.com/achebatarev/WallpaperSwitch/download"
	"github.com/achebatarev/WallpaperSwitch/setwall"
)

func HandleYesNoInput(s string) bool {
	var input string

	for true {
		fmt.Print(s)
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

func TUI() error {
	var wallpaper *download.Wallpaper

	done := false
	for !done {
		var err error

		wallpaper, err = download.DonwloadWallpaper()

		if err != nil {
			return fmt.Errorf("TUI: %w", err)
		}

		cmd, err := setwall.DisplayWallpaper(wallpaper)

		if err != nil {
			return fmt.Errorf("TUI: %w", err)
		}

		done = HandleYesNoInput("Do you want to set this image? (y/n) ")

		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("TUI: %w", err)
		}

	}

	if err := setwall.SwitchWallpaper(wallpaper); err != nil {
		return fmt.Errorf("TUI: %w", err)
	}

	if HandleYesNoInput("Do you want to add this image to favorites? (y/n) ") {
		if err := setwall.AddToFavorite(wallpaper); err != nil {
			return fmt.Errorf("TUI: %w", err)
		}
	}
	return nil
}
