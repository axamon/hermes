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

package hasher

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

// StringSum restituisce il chacksum hash della stringa passata come argomento.
func StringSum(str string) (hash string, err error) {

	h := md5.New()
	h.Write([]byte(str))

	hash = fmt.Sprintf("%x", h.Sum(nil))

	return hash, err
}

// StringSumWithSalt restituisce il chacksum hash della stringa passata come argomento
// più un seed personalizzato.
func StringSumWithSalt(str, salt string) (hash string, err error) {

	h := md5.New()
	h.Write([]byte(str + salt))

	hash = fmt.Sprintf("%x", h.Sum(nil))

	return hash, err
}

// FileSum restituisce il chacksum hash della stringa passata come argomento.
func FileSum(file string) (hash string, err error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	hash = fmt.Sprintf("%x", h.Sum(nil))

	return hash, err
}

// FileWithSeed restituisce il chacksum hash del file passato come argomento
// a cui viene aggiunto un seed.
func FileWithSeed(file, seed string) (hash string, err error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	hash = fmt.Sprintf("%x", h.Sum([]byte(seed)))

	return hash, err
}
