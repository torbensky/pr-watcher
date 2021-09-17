package prwatcher

import (
	"log"

	"github.com/gen2brain/beeep"
)

func Notify(title, description string) {
	err := beeep.Notify(title, description, "assets/OctoCat.png")
	if err != nil {
		log.Fatal("Error sending notification")
	}
}
