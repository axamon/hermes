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
			idv = "NON DISPOBINILE"
		}
		fmt.Println(idv)
}