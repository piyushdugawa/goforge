package utils

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
)

var destFolder = Gobin() // <- change as you like

const CONFIG_FILE = "GoForge.yaml"

var srcfilename = "./src/main.go"

var SrcContent = `package main

import "fmt"

func main(){
	fmt.Println("Project Initialised by GoForge!")
}`

var Pkg string

var cfgcontent = `  optimisation: true
  
  env:
    GOOS: windows
    GOARCH: amd64

  flags:
    - -ldflags
    - "-s -w"`

type Config struct {
	App struct {
		Package string `yaml:"package"`
		Version string `yaml:"version"`
	} `yaml:"app"`

	Build struct {
		Output       string `yaml:"output"`
		Optimisation bool   `yaml:"optimisation"`

		Env   map[string]interface{} `yaml:"env"`
		Flags []string               `yaml:"flags"`
	} `yaml:"build"`
}

func GetGOOSList(env map[string]interface{}) []string {
	if env == nil {
		return []string{}
	}
	val, ok := env["GOOS"]
	if !ok {
		return []string{}
	}
	switch v := val.(type) {
	case string:
		var list []string
		for _, s := range strings.Split(v, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				list = append(list, s)
			}
		}
		return list
	case []interface{}:
		var list []string
		for _, item := range v {
			s := strings.TrimSpace(fmt.Sprintf("%v", item))
			if s != "" {
				list = append(list, s)
			}
		}
		return list
	default:
		s := strings.TrimSpace(fmt.Sprintf("%v", val))
		if s != "" {
			return []string{s}
		}
		return []string{}
	}
}

func findHostPlatformInList(goosList []string) (string, int) {
	hostOS := runtime.GOOS
	for i, p := range goosList {
		if p == hostOS || (p == "mac" && hostOS == "darwin") {
			return p, i
		}
	}
	return "", -1
}

func compileTarget(cfg *Config, platform string, outputPath string) error {
	targetOS := platform
	if targetOS == "mac" {
		targetOS = "darwin"
	}

	originalEnv := make(map[string]string)
	for key, val := range cfg.Build.Env {
		if key == "GOOS" {
			continue
		}
		originalEnv[key] = os.Getenv(key)
		os.Setenv(key, fmt.Sprintf("%v", val))
	}
	if platform != "" {
		originalEnv["GOOS"] = os.Getenv("GOOS")
		os.Setenv("GOOS", targetOS)
	}

	defer func() {
		for key, val := range originalEnv {
			if val == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, val)
			}
		}
	}()

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(absOutput), 0755)
	if err != nil {
		return err
	}

	Chvenv("src")
	defer Chvenv("../")

	var args []string
	if cfg.Build.Optimisation {
		args = append(cfg.Build.Flags, "-o", absOutput, "main.go")
	} else {
		args = []string{"-o", absOutput, "main.go"}
	}

	color.Blue("🔨 Running build for platform %s:", platform)
	color.Cyan("go build %v\n", args)

	return Cmd("go", append([]string{"build"}, args...)...)
}

func GetTargetOutputPath(output string, platform string, isPrimary bool) string {
	if isPrimary {
		if platform == "windows" && !strings.HasSuffix(strings.ToLower(output), ".exe") {
			return output + ".exe"
		}
		return output
	}

	dir := filepath.Dir(output)
	base := filepath.Base(output)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	targetExt := ""
	if platform == "windows" {
		targetExt = ".exe"
	}
	return filepath.Join(dir, platform, name+targetExt)
}

func GetHostBinaryPath(cfg *Config) string {
	goosList := GetGOOSList(cfg.Build.Env)
	if len(goosList) == 0 {
		return GetTargetOutputPath(cfg.Build.Output, runtime.GOOS, true)
	}

	platform, index := findHostPlatformInList(goosList)
	if index == 0 {
		return GetTargetOutputPath(cfg.Build.Output, platform, true)
	} else if index > 0 {
		return GetTargetOutputPath(cfg.Build.Output, platform, false)
	} else {
		// Fallback: assume the primary target or host OS
		return GetTargetOutputPath(cfg.Build.Output, runtime.GOOS, true)
	}
}

func GenerateCfgContent(optimisation bool, selectedOSes []string) string {
	goosStr := fmt.Sprintf("[%s]", strings.Join(selectedOSes, ", "))

	return fmt.Sprintf(`  optimisation: %v
  
  env:
    GOOS: %s
    GOARCH: amd64

  flags:
    - -ldflags
    - "-s -w"`, optimisation, goosStr)
}

func PrioritizeHostOS(selectedOSes []string) []string {
	hostOS := runtime.GOOS
	var hostIdx = -1
	for i, osVal := range selectedOSes {
		if osVal == hostOS || (osVal == "mac" && hostOS == "darwin") {
			hostIdx = i
			break
		}
	}

	if hostIdx > 0 {
		hostVal := selectedOSes[hostIdx]
		selectedOSes = append(selectedOSes[:hostIdx], selectedOSes[hostIdx+1:]...)
		selectedOSes = append([]string{hostVal}, selectedOSes...)
	}

	return selectedOSes
}

func New(pkgName string) {
	var err error
	if pkgName == "" {
		pkgName, err = PromptPackageName()
		if err != nil {
			color.Red("❌ %v\n", err)
			return
		}
	}

	selectedOSes, err := PromptTargetOSes()
	if err != nil {
		color.Red("❌ Initialization cancelled.\n")
		return
	}

	selectedOSes = PrioritizeHostOS(selectedOSes)

	var optimisation bool
	err = huh.NewConfirm().
		Title("Enable build optimization?").
		Affirmative("Yes").
		Negative("No").
		Value(&optimisation).
		Run()
	if err != nil {
		color.Red("❌ Initialization cancelled.\n")
		return
	}

	customCfgContent := GenerateCfgContent(optimisation, selectedOSes)
	buildstr := filepath.Base(pkgName)
	preview := fmt.Sprintf("app:\n  package: %s\n  version: 0.0.1\nbuild:\n  output: build/%s\n%s", pkgName, buildstr, customCfgContent)

	confirmed := AskConfirm("Are you sure you want to initialise this module?", preview)
	if !confirmed {
		color.Red("❌ Initialization cancelled.\n")
		return
	}

	CreateCfgFile(CONFIG_FILE, pkgName, customCfgContent)

	dirname := "./src"

	// Check if directory exists
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		CreateFile(srcfilename, SrcContent)
	} else if err != nil {
		// Some other error occurred
		fmt.Println("Error checking directory:", err)
	} else if !info.IsDir() {
		// Path exists but is not a directory
		fmt.Printf("src is not a directory!\n")
	} else {
		// Directory exists, do nothing

		CreateFile(srcfilename, SrcContent)
	}

	Init()

}

func Build() {
	cfg, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		color.Red("❌ Failed to load config: %v\n", err)
		return
	}

	Chvenv("src")
	err = Cmd("go", "mod", "tidy")
	Chvenv("../")
	if err != nil {
		color.Red("❌ Error tidying modules: %v\n", err)
		return
	}

	goosList := GetGOOSList(cfg.Build.Env)

	// Check if this is an install request
	isInstall := len(os.Args) >= 2 && os.Args[1] == "install"

	if isInstall {
		hostOS := runtime.GOOS
		var outputPath string
		var platformToCompile string = hostOS

		if len(goosList) == 0 {
			outputPath = GetTargetOutputPath(cfg.Build.Output, hostOS, true)
		} else {
			platform, index := findHostPlatformInList(goosList)
			if index == 0 {
				outputPath = GetTargetOutputPath(cfg.Build.Output, platform, true)
				platformToCompile = platform
			} else if index > 0 {
				outputPath = GetTargetOutputPath(cfg.Build.Output, platform, false)
				platformToCompile = platform
			} else {
				outputPath = GetTargetOutputPath(cfg.Build.Output, hostOS, true)
			}
		}

		err = compileTarget(cfg, platformToCompile, outputPath)
		if err != nil {
			color.Red("❌ Build Failed for host platform %s: %v\n", hostOS, err)
			return
		}
		color.Green("✅ Host Build Successful!\n")
		return
	}

	if len(goosList) == 0 {
		err = compileTarget(cfg, "", GetTargetOutputPath(cfg.Build.Output, runtime.GOOS, true))
		if err != nil {
			color.Red("❌ Build Failed: %v\n", err)
			return
		}
		color.Green("✅ Build Successful!\n")
		return
	}

	for i, platform := range goosList {
		var outputPath string
		if i == 0 {
			outputPath = GetTargetOutputPath(cfg.Build.Output, platform, true)
		} else {
			outputPath = GetTargetOutputPath(cfg.Build.Output, platform, false)
		}

		err = compileTarget(cfg, platform, outputPath)
		if err != nil {
			color.Red("❌ Build Failed for platform %s: %v\n", platform, err)
			return
		}
	}
	color.Green("✅ All Builds Successful!\n")
}

func Run() {
	cfg, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		color.Red("❌ Failed to load config: %v\n", err)
		return
	}
	hostBin := GetHostBinaryPath(cfg)
	if !FileExists(hostBin) {
		Build()
		Run()
		return
	}
	color.Green("--------------------\n")
	color.Blue("Running Program...\n")
	color.Green("--------------------\n")
	fmt.Println()

	err = Cmd(hostBin)

	if err != nil {
		color.Red("Error running exe: %v\n", err)
		fmt.Println()
		return
	}
	fmt.Println()
}

func Install() {
	cfg, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		color.Red("❌ Failed to load config: %v\n", err)
		return
	}

	src := GetHostBinaryPath(cfg)

	// Resolve absolute paths
	absSrc, err := filepath.Abs(src)
	Check(err)

	// Verify source exists and is a regular file
	info, err := os.Stat(absSrc)
	Check(err)
	if info.IsDir() {
		fmt.Printf("Error: '%s' is a directory, not a file\n", absSrc)
		os.Exit(1)
	}

	// Ensure destination folder exists
	err = os.MkdirAll(destFolder, 0o755)
	Check(err)

	destPath := filepath.Join(destFolder, filepath.Base(absSrc))

	// Copy (overwrite if exists)
	err = CopyFile(absSrc, destPath)
	Check(err)

	fmt.Printf("Program installed to %s\n", destPath)
}

func Remove() {
	cfg, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		color.Red("❌ Failed to load config: %v\n", err)
		return
	}

	src := GetHostBinaryPath(cfg)

	binPath := Gobin()
	destPath := filepath.Join(binPath, filepath.Base(src))

	err = os.Remove(destPath)
	if err != nil {
		fmt.Println("Error removing program:", err)
		return
	}

	fmt.Println("Program removed successfully")
}

func Gobin() string {
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		// Fallback to GOPATH/bin if GOBIN is not set
		gopath := build.Default.GOPATH
		gobin = filepath.Join(gopath, "bin")
	}
	return gobin
}

func Help() {
	fmt.Printf(`Goforge - A minimal forge to build and manage your Go-based projects

Usage:
  goforge [command] [arguments]

Available Commands:
  help                 Show this help message
  version              Show the current version of goforge
  run                  Run the current project (main package)
  build                Build the project and output the executable
  new <pkg-name>       Initialize a new goforge project with the given package name
  install              Install project as a program in GOBIN
  remove               Remove the installed program from GOBIN
  clean	               Removes all builds and temporary files

For more information, visit: https://github.com/piyushdugawa/goforge/
`)
}

func Clean() {
	cfg, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		color.Red("❌ Failed to load config: %v\n", err)
		return
	}
	dir := filepath.Dir(cfg.Build.Output)
	if dir != "" && dir != "." && dir != "/" && dir != "\\" {
		err = os.RemoveAll(dir)
	} else {
		err = os.Remove(cfg.Build.Output)
		if err != nil && !os.IsNotExist(err) {
			// ignore
		}
		_ = os.Remove(cfg.Build.Output + ".exe")
	}
	if err != nil {
		fmt.Println("Error cleaning program:", err)
		return
	}

	fmt.Println("Project Cleaned Successfully!")
}
