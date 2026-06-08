package main

import (
	"GoForge/utils"
	"fmt"
	"os"

	"github.com/fatih/color"
)

const VERSION = "v0.10.1"
const VER_DESC = "restrict build/install to host OS during install"

func main() {
	if len(os.Args) < 2 {
		color.Yellow("⚠️  Usage: goforge <command> [arguments]\n")
		return
	}
	switch os.Args[1] {
	case "new":
		pkgName := ""
		if len(os.Args) >= 3 {
			pkgName = os.Args[2]
		}
		utils.New(pkgName)
	case "run":
		utils.Run()
	case "build":
		utils.Buildscr()
	case "install":
		utils.Buildscr()
		utils.Install()
	case "remove":
		utils.Remove()
	case "clean":
		utils.Clean()
	case "version":
		fmt.Println(VERSION)
	case "help":
		utils.Help()
	case "-h":
		utils.Help()
	case "--help":
		utils.Help()
	case "-v":
		fmt.Println(VERSION)
	case "--version":
		fmt.Println(VERSION)
	default:
		color.Yellow("⚠️  Usage: goforge <command> [arguments]\n")
		return
	}
}
