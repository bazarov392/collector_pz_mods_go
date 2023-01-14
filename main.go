package main

import (
	"golang-app/collector"
)

func main() {
	configData := collector.ConfigHandler()

	collect := collector.CollectorPzMods{
		ConfigData: configData,
		Params: collector.ResultParams{
			WorkshopIds: make([]string, 0),
			ModIds:      make([]string, 0),
		},
	}

	collect.GetWorkshopIds()
	collect.GetModIds()
	collect.GenerateFileModAndWorkshopIds()

}
