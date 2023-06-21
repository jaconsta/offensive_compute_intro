package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type demoRead struct{}

func (d *demoRead) Read(b []byte) (int, error) {
	fmt.Println("in> ")
	return os.Stdin.Read(b)
}

type demoWrite struct{}

func (d *demoWrite) Write(b []byte) (int, error) {
	fmt.Println("out> ")
	return os.Stdout.Write(b)
}

func main() {
}

func justIoCopy() {
	var (
		reader demoRead
		writer demoWrite
	)

	_, err := io.Copy(&writer, &reader)
	if err != nil {
		log.Fatalln("Unable to read/write data")
	}
}

func simpleImpl() {
	var (
		reader demoRead
		writer demoWrite
	)

	input := make([]byte, 4096)
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalln("Unable to read data")
	}
	fmt.Printf("Read %d bytes.\n", s)

	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("Unable to write data")
	}
	fmt.Printf("Wrote %d bytes", s)
}
