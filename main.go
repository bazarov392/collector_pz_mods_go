package main

import (
	"fmt"
	"golang-app/collector"
)

func main() {
	configData := collector.ConfigHandler()

	collect := collector.CollectorPzMods{
		ConfigData: configData,
		Params: collector.ResultParams{
			WorkshopIds: make([]string, 0),
			ModIds:      make([]string, 0),
			CopyInfo:    make([]collector.ModInfo, 0),
		},
	}

	collect.GetWorkshopIds()
	collect.GetModIds()
	collect.CopyModes()
	collect.GenerateFileModAndWorkshopIds()

	fmt.Println("Для закрытия нажмите ENTER")
	var end string
	fmt.Scanf("%s\n", &end)
}
