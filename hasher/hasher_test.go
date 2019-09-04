// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hasher_test

import (
	"fmt"
	"log"
	"os"

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

	file.Write([]byte("pippo"))

	file.Close()

	hash, err := hasher.FileWithSeed(testfile, "vvkidtbcjujhttuuikvjtfhilrkfkkfgejcktriignbr")
	if err != nil {
		log.Printf("ERROR Impossibile ricavare hash del file %s più seed: %s\n", testfile, err.Error())
	}

	fmt.Println(hash)

	err = os.Remove(testfile)
	if err != nil {
		log.Fatalf("ERROR impossibile cancellare file di test %s: %s\n", testfile, err.Error())
	}

	//Output:
	// 76766b69647462636a756a6874747575696b766a746668696c726b666b6b6667656a636b74726969676e6272a2242ead55c94c3deb7cf2340bfef9d5bcaca22dfe66e646745ee4371c633fc8

}
