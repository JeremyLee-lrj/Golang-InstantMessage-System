package main

import (
	"flag"
	"fmt"
)

func Testflag() {
	namePtr := flag.String("name", "Jeremy", "a string represents name")
	AgePtr := flag.Int("age", 23, "a integer represents the age")

	flag.Parse()
	fmt.Printf("name = %v, age = %v\n", *namePtr, *AgePtr)
}
