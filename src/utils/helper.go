package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func LoadConfig(path string) (*Config, error) { // return krte time config ko defreffrence krna h
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg) // unmarshalling me config ka ref aayega
	if err != nil {
		return nil, err
	}

	return &cfg, nil // ref return krna h
}

func Cmd(cmd string, str ...string) error {
	runcmd := exec.Command(cmd, str...)
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	return runcmd.Run()
}

func Chvenv(path string) {
	currPath, err := os.Getwd()
	if err != nil {
		color.Red("Error getting executable path:", err)
		return
	}

	// Construct absolute path to ./src/ relative to binary location
	virtualEnvPath := filepath.Join(currPath, path)

	// Check if it exists
	if _, err := os.Stat(virtualEnvPath); os.IsNotExist(err) {
		color.Red("Directory does not exist:", virtualEnvPath)
		return
	}

	// Change working directory
	err = os.Chdir(virtualEnvPath)
	if err != nil {
		color.Red("Error changing directory:", err)
		return
	}
	println(virtualEnvPath)
}

func CreateFile(srcfilename string, content string) {
	if _, err := os.Stat(srcfilename); os.IsNotExist(err) {
		// File does not exist, so create it
		file, err := os.Create(srcfilename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		_, err = file.WriteString(content)
		if err != nil {
			color.Red("❌ Error writing to file! %v\n", err)
			return
		}
	} else if err != nil {
		// Some other error occurred
		fmt.Println("Error checking file:", err)
	} else {
		fmt.Println("File already exists:", srcfilename)
	}
}

func CreateCfgFile(srcfilename string, pkgName string, srccontent string) {
	buildstr := path.Base(pkgName)
	content := "app:\n  package: " + pkgName + "\n" + "  version: 0.0.1" + "\n" + "build:\n  output: build/" + buildstr + "\n"
	if _, err := os.Stat(srcfilename); os.IsNotExist(err) {
		// File does not exist, so create it
		file, err := os.Create(srcfilename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		_, err = file.WriteString(content)
		if err != nil {
			color.Red("❌ Error writing to file! %v\n", err)
			return
		}
		_, err = file.WriteString(srccontent)
		if err != nil {
			color.Red("❌ Error writing to file! %v\n", err)
			return
		}
	} else if err != nil {
		// Some other error occurred
		fmt.Println("Error checking file:", err)
	} else {
		fmt.Println("File already exists:", srcfilename)
	}
}

func Init() {
	cfg, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		color.Red("❌ Error reading Config File: %v\n", err)
		return
	}
	Chvenv("src")
	err = Cmd("go", "mod", "init", cfg.App.Package)
	if err != nil {
		color.Red("❌ Error Initialising Project: %v\n", err)
		return
	}
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Force the same permissions as source (optional)
	info, err := os.Stat(src)
	if err == nil {
		_ = os.Chmod(dst, info.Mode())
	}
	return err
}

func Check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func Buildscr() {
	if len(os.Args) == 3 {
		switch os.Args[2] {
		case "run":
			Build()
			Chvenv("../")
			Run()
		default:
			color.Red("Argument '%v' not defined!\n", os.Args[2])
		}
	} else {
		Build()
	}
}
