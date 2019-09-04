// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/axamon/hermes/idvideoteca"
	"flag"
)

func main() {

	var rawurl = flag.String("u", "", "Url da gestire")

	flag.Parse()

		idv, err := idvideoteca.Find(*rawurl)
		if err != nil {
			fmt.Printf("NON DISPONIBILE: %s", err.Error())
			return
		}
		if idv == "" {
			fmt.Printf("NON DISPONIBILE: %s", *rawurl)
		}

		fmt.Println(idv)
}
