// Elvishd is a daemon for mediating access to the storage backend of elvish.
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/elves/elvish/elvishd/api"
	"github.com/elves/elvish/store"
)

var (
	elvishdPath  = flag.String("elvishd-path", "", "absolute path to the elvishd binary, required")
	socketPath   = flag.String("socket-path", "", "absolute path to the socket, required")
	databasePath = flag.String("database-path", "", "absolute path to the database, required")
	forked       = flag.Int("forked", 0, "how many times elvishd has forked, 0 (default), 1 or 2")
)

func flagsOK() bool {
	return 0 <= *forked && *forked <= 2 &&
		path.IsAbs(*elvishdPath) &&
		path.IsAbs(*socketPath) &&
		path.IsAbs(*databasePath)
}

const closeFd = ^uintptr(0)

func main() {
	flag.Parse()
	if !flagsOK() {
		flag.Usage()
		os.Exit(2)
	}
	switch *forked {
	case 0:
		syscall.Umask(0077)
		forkExecSelfAndExit(&syscall.ProcAttr{
			// cd to /
			Dir: "/",
			// empty environment
			Env: nil,
			// inherit stderr only
			Files: []uintptr{closeFd, closeFd, 2},
			Sys:   &syscall.SysProcAttr{Setsid: true},
		})
	case 1:
		forkExecSelfAndExit(nil)
	}
	serve(*socketPath, *databasePath)
}

func forkExecSelfAndExit(attr *syscall.ProcAttr) {
	_, err := syscall.ForkExec(*elvishdPath, []string{
		*elvishdPath,
		"-elvishd-path", *elvishdPath,
		"-socket-path", *socketPath,
		"-forked", strconv.Itoa(*forked + 1),
	}, attr)
	if err != nil {
		os.Stderr.WriteString("failed to ForkExec: " + err.Error() + "\n")
	}
	os.Exit(0)
}

func serve(socketPath string, databasePath string) {
	st, err := store.NewStore(databasePath)
	mustOK(err)

	listener, err := net.Listen("unix", socketPath)
	mustOK(err)
	defer os.Remove(socketPath)

	for {
		conn, err := listener.Accept()
		mustOK(err)
		go handle(conn, st)
	}
}

func handle(c net.Conn, st *store.Store) {
	defer c.Close()
	decoder := json.NewDecoder(c)
	encoder := json.NewEncoder(c)
	send := func(v interface{}) {
		err := encoder.Encode(v)
		if err != nil {
			log.Println("send:", err)
		}
	}
	for {
		var req api.Request
		err := decoder.Decode(&req)
		if err == io.EOF {
			return
		}
		if err != nil {
			send(&api.Response{Fatal: "decode: " + err.Error()})
			return
		}
		switch {
		case req.Ping:
			send(&api.Response{Number: 0})
		case req.ListDir:
			dirs, err := st.ListDirs()
			if err != nil {
				send(&api.Response{Error: "listdir: " + err.Error()})
				continue
			}
			send(&api.Response{Number: len(dirs)})
			for _, dir := range dirs {
				send(dir)
			}
		// case req.Quit:
		default:
			send(&api.Response{Error: "bad request"})
		}
	}
}

func mustOK(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
