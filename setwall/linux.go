package setwall

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/achebatarev/WallpaperSwitch/config"
	"github.com/achebatarev/WallpaperSwitch/download"
)

func SwitchWallpaper(wallpaper *download.Wallpaper) error {
	path := fmt.Sprint(config.Conf.Root, wallpaper.Name)
	arg0 := "--bg-scale"
	cmd := exec.Command("feh", arg0, path)

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("SwitchWallpaper %w", err)
	}

	return nil
}

func DisplayWallpaper(wallpaper *download.Wallpaper) (*exec.Cmd, error) {
	arg0 := "-x"
	cmd := exec.Command("feh", arg0, wallpaper.Preview)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("DisplayWallpaper %w", err)
	}

	return cmd, nil
}

func AddToFavorite(wallpaper *download.Wallpaper) error {

	srcPath := fmt.Sprint(config.Conf.Root, wallpaper.Name)
	dstPath := fmt.Sprint(config.Conf.Favorite, wallpaper.Name)

	srcFile, err := os.Open(srcPath)

	defer srcFile.Close()

	if err != nil {
		return fmt.Errorf("AddToFavorite: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	defer dstFile.Close()

	if err != nil {
		return fmt.Errorf("AddToFavorite: %w", err)
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("AddToFavorite: %w", err)
	}

	return nil
}
