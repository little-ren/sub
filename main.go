package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
var flagPrompt = flag.Bool("p", false, "prompt for each file")

func main() {

	// Parse parses the command-line flags from os.Args[1:]. Must be called
	// after all flags are defined and before flags are accessed by the program.
	flag.Parse()

	// string to []byte
	from := []byte(flag.Arg(0))

	// string to []byte
	to := []byte(flag.Arg(1))

	// *flagPrompt: dereference
	if err := s(from, to, *flagPrompt); err != nil {
		log.Fatal(err)
	}
}

func s(from []byte, to []byte, optPrompt bool) error {
	// Getwd returns a rooted path name corresponding to the
	// current directory.
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("s: error getting working directory: %v", err)
	}

	// Walk walks the file tree rooted at root, calling walkFn for each file or
	// directory in the tree, including root. All errors that arise visiting files
	// and directories are filtered by walkFn. The files are walked in lexical
	// order, which makes the output deterministic but means that for very
	// large directories Walk can be inefficient.
	// Walk does not follow symbolic links.
	if err := filepath.Walk(dir, sub(from, to, optPrompt)); err != nil {
		return fmt.Errorf("s: error walking filesytem %v", err)
	}

	return nil
}

func sub(from []byte, to []byte, optPrompt bool) filepath.WalkFunc {
	f := func(path string, info os.FileInfo, err error) error {
		n := info.Name()
		ignore := n == "" || n[0] == '.' || n[0] == '_' || n == "vendor"

		if info.IsDir() {
			if ignore {
				return filepath.SkipDir
			}
			return nil
		}

		if !ignore {
			if !optPrompt || prompt(path) {
				return subf(path, from, to)
			}
		}
		return nil
	}
	return filepath.WalkFunc(f)
}

func subf(path string, old []byte, new []byte) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("subf: error opening file: %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("subf: error reading file: %v", err)
	}

	bb := bytes.Replace(b, old, new, -1)

	if err := ioutil.WriteFile(path, bb, 0644); err != nil {
		return fmt.Errorf("subf: error writing file: %v", err)
	}
	return nil
}

func prompt(name string) bool {
	fmt.Printf("update file: %s (Y/N)", name)
	var confirm string
	fmt.Scanln(&confirm)
	switch confirm {
	case "y", "Y", "yes", "Yes", "YES":
		return true
	}
	return false
}
