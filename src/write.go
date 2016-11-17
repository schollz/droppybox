package gojot

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func WriteEntry() string {
	logger.Debug("Editing file")

	var cmdArgs []string
	if Editor == "vim" {
		// Setup vim
		vimrc := `set nocompatible
set backspace=2
func! WordProcessorModeCLI()
	setlocal formatoptions=t1
	setlocal textwidth=80
	map j gj
	map k gk
	set formatprg=par
	setlocal wrap
	setlocal linebreak
	setlocal noexpandtab
	normal G$
endfu
com! WPCLI call WordProcessorModeCLI()`
		// Append to .vimrc file
		if exists(path.Join(TempPath, ".vimrc")) {
			// Check if .vimrc file contains code
			logger.Debug("Found .vimrc.")
			fileContents, err := ioutil.ReadFile(path.Join(TempPath, ".vimrc"))
			if err != nil {
				log.Fatal(err)
			}
			if !strings.Contains(string(fileContents), "com! WPCLI call WordProcessorModeCLI") {
				// Append to fileContents
				logger.Debug("WPCLI not found in .vimrc, adding it...")
				newvimrc := string(fileContents) + "\n" + vimrc
				err := ioutil.WriteFile(path.Join(TempPath, ".vimrc"), []byte(newvimrc), 0644)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				logger.Debug("WPCLI found in .vimrc.")
			}
		} else {
			logger.Debug("Can not find .vimrc, creating new .vimrc...")
			err := ioutil.WriteFile(path.Join(TempPath, ".vimrc"), []byte(vimrc), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		cmdArgs = []string{"-u", path.Join(TempPath, ".vimrc"), "-c", "WPCLI", "+startinsert", path.Join(TempPath, "temp")}

	} else if Editor == "nano" {
		lines := "100" // TODO: DETERMINE THIS
		cmdArgs = []string{"+" + lines + ",1000000", "-r", "80", "--tempfile", path.Join(TempPath, "temp")}
	} else if Editor == "emacs" {
		lines := "100" // TODO: DETERMINE THIS
		cmdArgs = []string{"+" + lines + ":1000000", path.Join(TempPath, "temp")}
	} else if Editor == "micro" {
		settings := `{
    "autoclose": false,
    "autoindent": false,
    "colorscheme": "zenburn",
    "cursorline": false,
    "gofmt": false,
    "goimports": false,
    "ignorecase": false,
    "indentchar": " ",
    "linter": false,
    "ruler": false,
    "savecursor": false,
    "saveundo": false,
    "scrollmargin": 3,
    "scrollspeed": 2,
    "statusline": false,
    "syntax": false,
    "tabsize": 4,
    "tabstospaces": false,
		"softwrap": true
}`
		if !exists(path.Join(HomePath, ".config", "micro")) {
			os.MkdirAll(path.Join(HomePath, ".config", "micro"), 0755)
		}
		err := ioutil.WriteFile(path.Join(HomePath, ".config", "micro", "settings.json"), []byte(settings), 0644)
		if err != nil {
			log.Fatal(err)
		}

		lines := "10000000" // TODO determine this
		cmdArgs = []string{"-startpos", lines + ",1000000", path.Join(TempPath, "temp")}
	}

	// Load from binary assets
	logger.Debug("Trying to get asset: %s", "bin/"+Editor+Extension)
	data, err := Asset("bin/" + Editor + Extension)
	if err == nil {
		logger.Debug("Using builtin editor: %s", "bin/"+Editor+Extension)
		err = ioutil.WriteFile(path.Join(TempPath, Editor+Extension), data, 0755)
		if err != nil {
			log.Fatal(err)
		}
		Editor = path.Join(TempPath, Editor)
	} else {
		logger.Debug("Could not find builtin editor: %s", err.Error())
	}

	logger.Debug("Using editor %s", Editor)
	// Run the editor
	cmd := exec.Command(Editor+Extension, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		cmd := exec.Command(path.Join(ProgramPath, Editor+Extension), cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err2 := cmd.Run()
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	fileContents, _ := ioutil.ReadFile(path.Join(TempPath, "temp"))
	return strings.TrimSpace(string(fileContents))
}
