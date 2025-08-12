package main

import (
	"fmt"
	"os"

	go_m3u8 "github.com/globocom/go-m3u8"
)

func main() {
	file, _ := os.Open("multivariant.m3u8")

	fmt.Println("----- Decoding manifest to playlist -----")

	p, err := go_m3u8.ParsePlaylist(file)
	if err != nil {
		panic(err)
	}

	p.Print()

	fmt.Println("\n----- Encoding playlist back to string -----")

	m, err := go_m3u8.EncodePlaylist(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(m)
}
