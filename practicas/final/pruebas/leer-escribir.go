package main

import (
    "io"
    "os"
    // "fmt"
)

func main() {
    infile := "input.txt"
    outfile := "output.txt"
    if _, err := os.Stat("./"+infile); os.IsNotExist(err) {
        os.Create(infile)
    }
    fi, err := os.Open(infile)
    if err != nil {
        panic(err)
    }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()

    // open output file
    fo, err := os.Create(outfile)
    if err != nil {
        panic(err)
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()

    // make a buffer to keep chunks that are read
    buf := make([]byte, 1024)
    for {
        // read a chunk
        n, err := fi.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if n == 0 {
            break
        }

        // write a chunk
        if _, err := fo.Write(buf[:n]); err != nil {
            panic(err)
        }
    }
}
