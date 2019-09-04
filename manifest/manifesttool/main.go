// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"fmt"
	"github.com/axamon/hermes/manifest"
)

func main() {

var rawurl = flag.String("u", "", "Url da gestire")

flag.Parse()
	
manifesturl, err := manifest.Find(*rawurl)
if err != nil {
	log.Println(err.Error())
}

fmt.Println(manifesturl)
}