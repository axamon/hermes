// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hasher_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/axamon/hermes/hasher"
)

func ExampleFileSum() {
	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}

	file.Write([]byte("pippo"))

	file.Close()

	hash, err := hasher.FileSum(testfile)
	if err != nil {
		log.Printf("ERROR Impossibile ricavare hash del file %s: %s\n", testfile, err.Error())
	}

	fmt.Println(hash)

	err = os.Remove(testfile)
	if err != nil {
		log.Fatalf("ERROR impossibile cancellare file di test %s: %s\n", testfile, err.Error())
	}
	//Output:
	// a2242ead55c94c3deb7cf2340bfef9d5bcaca22dfe66e646745ee4371c633fc8
}

func ExampleFileWithSeed() {
	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}
	defer func() {
		err = os.Remove(testfile)

		if err != nil {
			log.Fatalf("ERROR impossibile cancellare file di test %s: %s\n", testfile, err.Error())
		}
	}()

	file.Write([]byte("pippo"))

	file.Close()

	hash, err := hasher.FileWithSeed(testfile, "vvkidtbcjujhttuuikvjtfhilrkfkkfgejcktriignbr")
	if err != nil {
		log.Printf("ERROR Impossibile ricavare hash del file %s più seed: %s\n", testfile, err.Error())
	}

	fmt.Println(hash)

	//Output:
	// 76766b69647462636a756a6874747575696b766a746668696c726b666b6b6667656a636b74726969676e6272a2242ead55c94c3deb7cf2340bfef9d5bcaca22dfe66e646745ee4371c633fc8

}

func TestFileSum(t *testing.T) {

	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}
	defer func() {
		err = os.Remove(testfile)

		if err != nil {
			log.Fatalf("ERROR impossibile cancellare file di test %s: %s\n", testfile, err.Error())
		}
	}()

	file.Write([]byte("pippo"))

	file.Close()

	type args struct {
		file string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		{"Primo", args{file: "test.zip"}, "a2242ead55c94c3deb7cf2340bfef9d5bcaca22dfe66e646745ee4371c633fc8", false},
		{"Primo", args{file: "test1.zip"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := hasher.FileSum(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("FileSum() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestStringSum(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		{"Primo", args{str: "pippo"}, "0c88028bf3aa6a6a143ed846f2be1ea4", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := hasher.StringSum(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("StringSum() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestStringSumWithSalt(t *testing.T) {
	type args struct {
		str  string
		salt string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		{"Primo", args{str: "pippo", salt: "Tim1"}, "086b4291bed0b78adbde667f921b6d4a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := hasher.StringSumWithSalt(tt.args.str, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSumWithSalt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("StringSumWithSalt() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
