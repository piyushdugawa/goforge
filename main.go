package main

import (
	"github.com/piyushdugawa/girocco/cmd"
)

const VERSION = "v0.14.0"
const VER_DESC = "Project name changed from GoForge to Girocco due to some conflicts."

func main() {
	cmd.Execute(VERSION)
}
