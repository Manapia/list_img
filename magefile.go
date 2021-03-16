// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	packageName = "github.com/Manapia/list_img"
	binFileName = "list_img"
)

func Build() error {
	fmt.Println("Building...")
	binName := genRelBinName(runtime.GOOS, runtime.GOARCH)
	return buildBase(binName)
}

func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	binName := genRelBinName(runtime.GOOS, runtime.GOARCH)
	_, filename := filepath.Split(binName)

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	return sh.Copy(filepath.Join(goPath, "bin", filename), binName)
}

func BuildForWindows() error {
	fmt.Println("Building for windows...")
	env := map[string]string{"GOOS": "windows", "GOARCH": "amd64"}
	binName := genRelBinName(env["GOOS"], env["GOARCH"])
	return envBuildBase(env, binName)
}

func BuildForDarwin() error {
	fmt.Println("Building for darwin...")
	env := map[string]string{"GOOS": "darwin", "GOARCH": "amd64"}
	binName := genRelBinName(env["GOOS"], env["GOARCH"])
	return envBuildBase(env, binName)
}

func BuildForLinux() error {
	fmt.Println("Building for linux...")
	env := map[string]string{"GOOS": "linux", "GOARCH": "amd64"}
	binName := genRelBinName(env["GOOS"], env["GOARCH"])
	return envBuildBase(env, binName)
}

func Clean() error {
	fmt.Println("Cleaning...")
	return os.RemoveAll("build")
}

func buildBase(binName string) error {
	return sh.Run("go", "build", "-ldflags", "-s -w", "-o", binName, packageName)
}

func envBuildBase(env map[string]string, binName string) error {
	return sh.RunWith(env, "go", "build", "-ldflags", "-s -w", "-o", binName, packageName)
}

func genRelBinName(goOs, goArch string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("cannot get current directory:", err)
	}

	relPath := filepath.Join(currentDir, "build", goOs+"_"+goArch, binFileName)
	if goOs == "windows" {
		relPath += ".exe"
	}
	return relPath
}
