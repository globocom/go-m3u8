package main

import (
	"os"

	go_m3u8 "github.com/globocom/go-m3u8"
)

func main() {
	file, _ := os.Open("multivariant.m3u8")
	p, err := go_m3u8.ParsePlaylist(file)

	if err != nil {
		panic(err)
	}

	p.Print()
}
