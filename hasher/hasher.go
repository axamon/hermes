// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
		// log.Fatal(err)
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		// log.Fatal(err)
		return "", err
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
