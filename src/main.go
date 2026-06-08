package main

import (
	"GoForge/cmd"
)

const VERSION = "v0.11.0"
const VER_DESC = "restrict build/install to host OS during install"

func main() {
	cmd.Execute(VERSION)
}
