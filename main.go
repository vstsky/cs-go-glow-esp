package main

import (
	"csgo/cs"
	"fmt"
	"time"
)

var appTitle = "CSGO"

func main() {
	err := cs.Init()

	if err != nil {
		fmt.Println("couldn't open process")
	}

	player := cs.GetLocalPlayer()

	var enemies []cs.Entity

	for {
		enemies = []cs.Entity{}

		for _, entity := range cs.GetEntities() {
			if entity.Team != player.Team && entity.Health > 0 {
				enemies = append(enemies, entity)
			}
		}

		for _, enemy := range enemies {
			enemy.EnableGlowEsp()
		}

		time.Sleep(time.Millisecond * 2)
	}
}
