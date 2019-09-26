package main

import "log"

func main() {
	config := ReadConfig()
	torrents := GetTorrents(&config.Transmission)

	for index, name := range torrents {
		log.Printf("%d %s", index, name)
	}
	Consume(&config.Sqs)
}
