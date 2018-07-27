package main

import (
	"event_processor/plugin"
	"log"
	"time"

	"event_processor"
)

var plugins plugin.Plugins

func main() {
	plugins = make([]event_processor.Plugin, 0)
	plugins.LoadPlugins()
	plugins.SetupPlugins()

	log.Printf("%+v\n", plugins)

	for {
		log.Println("Hello Fresh")
		for _, plugin := range plugins {
			plugin.Process(nil)
		}
		time.Sleep(15 * time.Second)
	}
}
