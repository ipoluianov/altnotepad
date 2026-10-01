package app

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"time"
)

// A second start of the application passes its files to the running one
// through a local socket, then exits (like Notepad++)

func socketPath(dir string) string {
	return filepath.Join(dir, "instance.sock")
}

// OpenRequest is what a second start asks the running one to open
type OpenRequest struct {
	Files []string
	Line  int `json:",omitempty"`
}

// SendToRunning passes the request to the running application; false if none runs
func SendToRunning(dir string, req OpenRequest) bool {
	conn, err := net.DialTimeout("unix", socketPath(dir), 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	bs, _ := json.Marshal(req)
	if _, err := conn.Write(append(bs, '\n')); err != nil {
		return false
	}
	// Wait for the answer, so the request is not lost if the running one is closing
	line, err := bufio.NewReader(conn).ReadString('\n')
	return err == nil && line == "ok\n"
}

// ListenForInstances receives the requests of the later starts; onOpen runs on the listening goroutine
func ListenForInstances(dir string, onOpen func(req OpenRequest)) func() {
	os.MkdirAll(dir, 0755)
	path := socketPath(dir)
	os.Remove(path) // left by a crash
	ln, err := net.Listen("unix", path)
	if err != nil {
		return func() {}
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				conn.SetDeadline(time.Now().Add(5 * time.Second))
				line, err := bufio.NewReader(conn).ReadBytes('\n')
				if err != nil {
					return
				}
				var req OpenRequest
				if json.Unmarshal(line, &req) != nil {
					return
				}
				conn.Write([]byte("ok\n"))
				onOpen(req)
			}(conn)
		}
	}()
	return func() {
		ln.Close()
		os.Remove(path)
	}
}
