package main

import (
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
)

func main() {
	// create the client
	client, err := torrent.NewClient(torrent.NewDefaultClientConfig())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// scan for the file name
	fmt.Printf("Enter magnet link: ")
	var magnetLink string
	fmt.Scanln(&magnetLink)

	// add torrent magnet
	t, err := client.AddMagnet(magnetLink)
	if err != nil {
		panic(err)
	}

	fmt.Println("Fetching metadata...")

	// Wait for the metadata to be ready
	select {
	case <-t.GotInfo():
		fmt.Println("Metadata fetched successfully.")
	case <-time.After(30 * time.Second):
		fmt.Println("Timeout while fetching metadata.")
		return
	}

	// download torrent file
	t.DownloadAll()

	// show progress
	totalSize := t.Info().TotalLength()
	var lastBytesCompleted int64 = 0
	for {
		bytesCompleted := t.BytesCompleted()
		progress := float64(bytesCompleted) / float64(totalSize) * 100
		speed := float64(bytesCompleted-lastBytesCompleted) / (1024 * 1024) // MB/s
		lastBytesCompleted = bytesCompleted
		fmt.Printf("\rProgress: %.2f%% | Speed: %.2f MB/s", progress, speed)
		if bytesCompleted >= totalSize { // Check if download is complete
			fmt.Println("\nDownload complete.")
			break
		}
		time.Sleep(1 * time.Second) // Avoid spamming the output
	}
}
