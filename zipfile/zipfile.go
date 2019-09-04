// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zipfile

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// ReadAll legge il file zippato passato come parametro e restituisce
// un io.Reader e un eventuale errore.
func ReadAll(ctx context.Context, zipFile string) (content io.Reader, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	// Apre il file zippato in lettura.
	f, err := os.Open(zipFile)
	defer f.Close()
	if err != nil {
		log.Printf("Error Impossibile aprire il file %s: %s\n", zipFile, err.Error())
	}

	// Unzippa in memoria.
	gr, err := gzip.NewReader(f)
	defer gr.Close()
	if err != nil {
		log.Printf("Error Impossibile leggere il contenuto del file %s: %s\n", zipFile, err.Error())
	}

	return gr, err
}

// ReadAllGZ legge il file compresso in gzip passato come parametro
// e restituisce l'intero contenuto del file e un eventuale errore.
func ReadAllGZ(ctx context.Context, zipFile string) (content []byte, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	// Apre il file zippato in lettura.
	f, err := os.Open(zipFile)
	defer f.Close()
	if err != nil {
		log.Printf("Error Impossibile aprire il file %s: %s\n", zipFile, err.Error())
	}

	// Unzippa in memoria.
	gr, err := gzip.NewReader(f)
	defer gr.Close()
	if err != nil {
		log.Printf("Error Impossibile leggere il contenuto del file %s: %s\n", zipFile, err.Error())
	}

	content, err = ioutil.ReadAll(gr)
	if err != nil {
		log.Printf("Error Impossiblile copiare in memoria il file: %s\n", err.Error())
	}

	return content, err
}

// ReadAllZIP legge il file compresso in ZIP passato come parametro e restituisce
// l'intero contenuto del file e un eventuale errore.
// func ReadAllZIP(ctx context.Context, zipFile string) (content []byte, err error) {

// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("Recovered in f", r)
// 		}
// 	}()

// 	// Unzippa in memoria.
// 	r, err := zip.OpenReader(zipFile)
// 	defer r.Close()
// 	if err != nil {
// 		log.Printf("Error Impossibile leggere il contenuto del file %s: %s\n", zipFile, err.Error())
// 	}

// 	for _, f := range r.File {
// 		// Store filename/path for returning and using later on
// 		if strings.Contains(f.Name, "csv") {
// 			rc, err := f.Open()
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			content, err = ioutil.ReadAll(rc)
// 			if err != nil {
// 				log.Printf("Error Impossiblile copiare in memoria il file: %s\n", err.Error())
// 			}
// 			break
// 		}
// 	}

// 	return content, err
// }
