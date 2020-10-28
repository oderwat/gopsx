// 2>/dev/null;/usr/bin/env test -x $0 && (go build -o "${0}c" "$0" && "${0}c" "$@"; r=$?; rm -f "${0}c"; exit "$r"); exit "$?"
//
// The Shebang is inspired by:
// https://gist.github.com/eSlider/b8cd54ab25600ce9e8098b7fe55a9e29
//
// This little helper is a wrapper for gops to make it easier
// to use by allowing the pid to be replaces with a glob pattern process name

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	switch len(os.Args) {
	case 1:
		execLive([]string{"gops"}...)
		os.Exit(0)
	case 2:
		switch os.Args[1] {
		case "help":
			fallthrough
		case "tree":
			execLive([]string{"gops", os.Args[1]}...)
			os.Exit(0)
		}
	}

	// we have at least two args so we search for the last one as process

	// list all of them into a buffer
	var cmd = execCapture([]string{"gops"}...)

	scanner := bufio.NewScanner(cmd.Stdout.(*bytes.Buffer))

	type nameAndPid struct {
		pid   string
		name  string
		query bool
	}

	var pids []nameAndPid

	searchNr := 0

	lookup := os.Args[len(os.Args)-1]
	// check if we actually got a pid!
	if _, err := strconv.Atoi(lookup); err == nil {
		pids = append(pids, nameAndPid{lookup, lookup, true})
	} else {
		onlyAgent := false
		searchPat := ""
		search := strings.SplitN(lookup, "@", 2)
		searchPat = search[0]
		if len(search) == 2 {
			if search[1] == "*" {
				searchNr = -1
			} else if search[1] == "+" {
				searchNr = -1
				onlyAgent = true
			} else {
				nr, err := strconv.Atoi(search[1])
				if err != nil {
					log.Fatalln("Parameter after @ must be an int")
				}
				searchNr = nr
			}
		}

		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			query := fields[3] == "*"
			if onlyAgent && !query {
				continue
			}

			// we skip gops because it will vanish!
			if fields[2] == "gops" {
				continue
			}

			match, err := filepath.Match(searchPat, fields[2])
			if err != nil {
				log.Fatalf("filepath.Match() failed with %s\n", err)
			}
			if match {
				pids = append(pids, nameAndPid{fields[0], fields[2], query})
				if searchNr > 0 && len(pids) == searchNr {
					// no more needed
					break
				}
			}
		}
		if len(pids) == 0 {
			fmt.Println("No matching Go process was found")
			os.Exit(1)
		}

		if searchNr == 0 && len(pids) > 1 {
			fmt.Println("Multiple matches where found!\nUse name@# (for position #) or name@* (for all) or the PID")
			for i, nap := range pids {
				fmt.Printf("%d\t%s\t%s\n", i+1, nap.pid, nap.name)
			}
			os.Exit(1)
		}

		if searchNr != -1 && searchNr-1 > len(pids) {
			fmt.Printf("A match with offset %d was not found\n", searchNr)
			for i, nap := range pids {
				fmt.Printf("%d\t%s\t%s\n", i, nap.pid, nap.name)
			}
			os.Exit(1)
		}

		if searchNr != 0 && searchNr != -1 {
			pids = []nameAndPid{pids[searchNr-1]}
		}
	}

	if os.Args[1] == "pid" {
		if len(pids) > 1 {
			for _, p := range pids {
				fmt.Println(p.pid)
			}
		} else {
			fmt.Print(pids[0].pid)
		}
		return
	}

	os.Args[0] = "gops"
	for i, p := range pids {
		if i > 0 {
			fmt.Print("\n")
		}
		if len(pids) > 1 || searchNr == -1 {
			// if we have multiple results we print a header
			fmt.Printf("pid: %s (%s)\n", p.pid, p.name)
		}
		os.Args[len(os.Args)-1] = p.pid
		if len(os.Args) > 2 && !p.query {
			fmt.Println("(no agent)")
			continue
		}
		execLive(os.Args...)
	}
	os.Exit(0)
}

// execCapture runs a command with some arguments and captures stdout and stderr
func execCapture(arg ...string) *exec.Cmd {

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(arg[0], arg[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return cmd
}

// execLive runs a command with some arguments directly to stdout / stderr
func execLive(arg ...string) {
	cmd := exec.Command(arg[0], arg[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
