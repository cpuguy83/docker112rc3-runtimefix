package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type skipped struct {
	id  string
	msg string
}

func (s *skipped) Error() string {
	return fmt.Sprintf("skipping container '%s': %s", s.id, s.msg)
}

func main() {
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if len(dir) == 0 {
		noDir()
	}

	dir = filepath.Join(dir, "containers")
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading from containers dir: %v\n", err)
		os.Exit(3)
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "could not change dirs: %v\n", err)
		os.Exit(2)
	}

	var numProcessed int
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		if err := process(info); err != nil {
			if _, ok := err.(*skipped); ok {
				fmt.Fprintln(os.Stdout, err)
			} else {
				fmt.Fprintln(os.Stderr, "error while patching container config:", err)
			}
			continue
		}
		fmt.Fprintln(os.Stdout, "proccessed", info.Name())
		numProcessed++
	}

	c := "containers"
	if numProcessed == 1 {
		c = "container"
	}
	fmt.Fprintf(os.Stderr, "Processed %d %s\n", numProcessed, c)
	if numProcessed > 0 {
		fmt.Fprintln(os.Stderr, "You must now restart docker for changes to take affect")
	}
}

func noDir() {
	fmt.Fprintln(os.Stderr, "Must provide 1 argument specifying the path to the containers directory")
	os.Exit(1)
}

func process(info os.FileInfo) error {
	path := info.Name()
	id := filepath.Base(path)
	hcPath := filepath.Join(path, "hostconfig.json")

	b, err := ioutil.ReadFile(hcPath)
	if err != nil {
		return fmt.Errorf("could not find config for container '%s'\n: %v", id, err)
	}

	var hc HostConfig
	if err := json.Unmarshal(b, &hc); err != nil {
		return fmt.Errorf("could not read config for container %s: %v\n", id, err)
	}

	if hc.Runtime != "default" {
		return &skipped{id, "already fixed"}
	}

	hc.Runtime = "runc"
	b, err = json.Marshal(&hc)
	if err != nil {
		return fmt.Errorf("could not fix config for container %s: %v\n", id, err)
	}
	hcStat, err := os.Stat(hcPath)
	if err != nil {
		return fmt.Errorf("could not read file mode for container config '%s', skipping: %v\n", id, err)
	}
	if err := AtomicWriteFile(hcPath, b, hcStat.Mode()); err != nil {
		return fmt.Errorf("could not fix config for container %s: %v\n", id, err)
	}

	// double-check and make sure we didn't fubar something
	if b, _ := ioutil.ReadFile(hcPath); b != nil {
		var hc HostConfig
		if err := json.Unmarshal(b, &hc); err == nil {
			if hc.Runtime != "runc" {
				return fmt.Errorf("unexpected value from 'Runtime' field for container '%s', config:\n", id, string(b))
			}
		}
	}
	return nil
}
