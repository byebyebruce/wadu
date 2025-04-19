package main

import (
	"fmt"

	"github.com/byebyebruce/wadu/tts/edgetts"
)

func main() {
	// List all available voices
	voices, err := edgetts.ListVoices()
	if err != nil {
		panic(err)
	}

	// Print the voices
	for _, voice := range voices {
		fmt.Println(voice)
	}
}
