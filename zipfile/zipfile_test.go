package zipfile_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/zipfile"
)

const data = `
test
test
test
`

func ExampleReadAll() {
	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}

	zippato := gzip.NewWriter(file)

	_, err = zippato.Write([]byte(data))
	zippato.Close()

	content, err := zipfile.ReadAll(testfile)

	scan := bufio.NewScanner(content)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println(line)
	}
	file.Close()
	err = os.Remove(testfile)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}
	// Output:
	// test
	// test
	// test
}
