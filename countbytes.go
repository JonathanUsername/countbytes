package main

// This package simply takes a streaming input from stdin and quietly 
// passes it through to stdout while outputting the total bytes so far 
// transferred as stderr. I wrote this to unobtrusively monitor 
// transferal of large files via, e.g. rsync 

import (
    "bufio"
    "os"
    "os/signal"
    "syscall"
    "io"
    "log"
    "fmt"
    "github.com/dustin/go-humanize"
)

// PassThru wraps an existing io.Reader.
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PassThru struct {
    io.Reader
    total uint64 // Total # of bytes transferred
}

// Read 'overrides' the underlying io.Reader's Read method.
 // We simply use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
  n, err := pt.Reader.Read(p)
  pt.total += uint64(n)

  if err == nil {
    output := fmt.Sprintf("(%s) ", humanize.Bytes(pt.total))
    fmt.Fprintf(os.Stderr, "\x1B[2K\x1B[%dD%s", len(output), output)
  }

  return n, err
}

func captureExit () {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func(){
    <-c
    fmt.Println("") // Print empty line to avoid ANSI oddities
    os.Exit(1)
  }()
}

func main () {
  var src io.Reader
  reader := bufio.NewReader(os.Stdin)
  src = &PassThru{Reader: reader}
  captureExit()
  for {
    buf := make([]byte, 0, 4*1024) // Empty byte array of len 0, cap 4096
    v, err := src.Read(buf[:cap(buf)])
    buf = buf[:v]
    if v == 0 {
      if err == nil {
        continue
      }
      if err == io.EOF {
          break
      }
      log.Fatal(err)
    }
    os.Stdout.Write(buf)
  }
}

