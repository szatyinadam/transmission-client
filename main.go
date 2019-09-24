package main

import "log"

func main() {
	torrents := GetTorrents()

	for index, name := range torrents {
		log.Printf("%d %s", index, name)
	}
}
