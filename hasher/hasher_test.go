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

package hasher_test

import (
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/hasher"
)

func ExampleSum() {
	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}

	file.Write([]byte("pippo"))

	file.Close()

	hash, err := hasher.Sum(testfile)
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

func ExampleWithSeed() {
	testfile := "test.zip"
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(testfile, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}

	file.Write([]byte("pippo"))

	file.Close()

	hash, err := hasher.WithSeed(testfile, "vvkidtbcjujhttuuikvjtfhilrkfkkfgejcktriignbr")
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
