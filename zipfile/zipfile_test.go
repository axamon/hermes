// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package zipfile_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

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

func TestReadAll(t *testing.T) {
	type args struct {
		zipFile string
	}
	tests := []struct {
		name        string
		args        args
		wantContent io.Reader
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := ReadAll(tt.args.zipFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("ReadAll() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}
