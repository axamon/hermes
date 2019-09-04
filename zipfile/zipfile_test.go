// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zipfile_test

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
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

func ExampleReadAllGZ() {
	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}

	zippato := gzip.NewWriter(file)

	_, err = zippato.Write([]byte(data))
	zippato.Close()

	ctx := context.TODO()
	content, err := zipfile.ReadAll(ctx, testfile)

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

func TestReadAllGZ(t *testing.T) {

	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}

	zippato := gzip.NewWriter(file)

	_, err = zippato.Write([]byte(data))
	zippato.Close()

	file.Close()

	defer func() {
		err = os.Remove(testfile)
		if err != nil {
			log.Fatalf("Failed to open zip for writing: %s", err)
		}
	}()

	type args struct {
		ctx     context.Context
		zipFile string
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
		wantErr     bool
	}{
		{"Primo", args{ctx: context.TODO(), zipFile: "test.zip"}, []byte(data), false},
		{"fileinesistente", args{ctx: context.TODO(), zipFile: "test1.zip"}, []byte(nil), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := zipfile.ReadAllGZ(tt.args.ctx, tt.args.zipFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAllGZ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("ReadAllGZ() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}
