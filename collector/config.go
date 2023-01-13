package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const PathToConfigFile = "/tmp/collectorPzModsForMautConfig"

type Config struct {
	SteamDirMods         string
	ResultDirMods        string
	SaveDataMode         string
	PathServerIniFile    string
	IgnoreModIds         []string
	IgnoreWorkshopIds    []string
	CopyFiles            bool
	GenerateModsList     bool
	GenerateWorkshopList bool
}

func ConfigHandler() Config {

	configData := Config{
		SteamDirMods:         "",
		ResultDirMods:        "",
		SaveDataMode:         "",
		PathServerIniFile:    "",
		IgnoreWorkshopIds:    make([]string, 0),
		IgnoreModIds:         make([]string, 0),
		CopyFiles:            true,
		GenerateModsList:     true,
		GenerateWorkshopList: true,
	}
	if isConfigFile() {
		configFile, err := ioutil.ReadFile(PathToConfigFile)
		if err != nil {
			os.Remove(PathToConfigFile)
			panic("Не удалось прочитать конфиг файл, запустите программу снова")
		}

		err2 := json.Unmarshal(configFile, &configData)
		if err2 != nil {
			os.Remove(PathToConfigFile)
			panic("Не удалось прочитать конфиг файл, запустите программу снова")
		}

		var isGenerateNewConfig string

		fmt.Println("Использовать ли сохраненные настройки? (Д - если использовать, Н если задать другое)")
		fmt.Scanf("%s\n", &isGenerateNewConfig)

		if !strings.Contains(isGenerateNewConfig, "Д") {
			generateConfig(&configData)
		}
	} else {
		generateConfig(&configData)
	}

	return configData
}

func isConfigFile() bool {
	_, err := os.Open(PathToConfigFile)
	return err == nil
}

func generateConfig(configData *Config) {
	fmt.Println("Настройки конструктора не были найдены, необходимо их ввести для продолжения.")

	fmt.Println("Путь до папки с модами Project Zomboid:")
	fmt.Scanf("%s\n", &configData.SteamDirMods)

	fmt.Println("Путь до папки, куда необходимо всё выложить: ")
	fmt.Scanf("%s\n", &configData.ResultDirMods)

	fmt.Println("Перечислите ID модов, которые необходимо проигнорировать, через \",\": ")
	var stringIgnoreModIds string
	fmt.Scanf("%s\n", &stringIgnoreModIds)
	configData.IgnoreModIds = strings.Split(stringIgnoreModIds, ",")

	fmt.Println("Перечислите Workshop ID модов, которые необходимо проигнорировать, через \",\": ")
	var stringIgnoreWorkshopIds string
	fmt.Scanf("%s\n", &stringIgnoreWorkshopIds)
	configData.IgnoreWorkshopIds = strings.Split(stringIgnoreWorkshopIds, ",")

	configData.GenerateWorkshopList = getTrueOrFalseFromConsole("Сохранить Workshop ID в отдельном файле? (Н - Нет, Д - Да)")
	configData.GenerateModsList = getTrueOrFalseFromConsole("Сохранить Mod ID в отдельном файле? (Н - Нет, Д - Да)")
	configData.CopyFiles = getTrueOrFalseFromConsole("Собирать ли моды? (Н - Нет, Д - Да)")

	jsonData, err := json.Marshal(configData)
	if err != nil {
		panic(1)
	}

	configFile, err := os.Create(PathToConfigFile)
	if err != nil {
		panic("Не удалось создать файл с настройками, запустите скрипт с вставкой \"sudo\"")
	}

	configFile.WriteString(string(jsonData))

	configFile.Close()
	fmt.Println("Настройки сохранены.")
}

func getTrueOrFalseFromConsole(message string) bool {
	var result string
	fmt.Scanf("%s\n", &result)

	if result == "Д" {
		return true
	} else if result == "Н" {
		return false
	} else {
		fmt.Printf("Не удалось прочитать ответ, он должен быть либо \"Д\", либо \"Н\"")
		return getTrueOrFalseFromConsole(message)
	}

}
