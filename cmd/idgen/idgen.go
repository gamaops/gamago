package main

import (
	"fmt"

	"flag"

	"github.com/gamaops/gamago/pkg/id"
)

func main() {

	outputEncoding := flag.String("encoding", "base32", "Encoding to output the generated IDs")
	count := flag.Int("count", 2, "How many IDs you want")
	randomSize := flag.Int("random", 5, "Random part size")

	flag.Parse()

	gen, err := id.NewIDGenerator(int32(*randomSize))
	if err != nil {
		panic(err)
	}
	for i := 0; i < *count; i++ {

		id, err := gen.New()

		if err != nil {
			panic(err)
		}

		var idStr string

		if *outputEncoding == "base32" {
			idStr = id.Base32()
		} else {
			idStr = id.Hex()
		}

		fmt.Printf("%v\n", idStr)

	}
}
