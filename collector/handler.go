package collector

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type CollectorPzMods struct {
	ConfigData Config
	Params     ResultParams
}

type ResultParams struct {
	WorkshopIds []string
	ModIds      []string
}

func (c *CollectorPzMods) GetWorkshopIds() {
	dirs, err := ioutil.ReadDir(c.ConfigData.SteamDirMods)
	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			c.Params.WorkshopIds = append(c.Params.WorkshopIds, dir.Name())
			fmt.Printf("Workshop ID %s добавлен в обработчик\n", dir.Name())
		}

	}

}

func (c *CollectorPzMods) GetModIds() {
	ch := make(chan string)
	for _, workshopDir := range c.Params.WorkshopIds {
		pathWorkshopDir := c.ConfigData.SteamDirMods + "/" + workshopDir

		mods, err := ioutil.ReadDir(pathWorkshopDir + "/mods")
		if err != nil {
			panic("Неудалось прочесть директорию " + pathWorkshopDir + "/mods")
		}

		for _, modDir := range mods {

			pathDirMode := pathWorkshopDir + "/mods/" + modDir.Name()
			contentModInfo, err := ioutil.ReadFile(pathDirMode + "/mod.info")
			if err != nil {
				panic("Не удалось прочитать файл: " + pathDirMode + "/mod.info")
			}

			re, _ := regexp.Compile("id=(.+)(\r\n|\n)?")
			stringModId := strings.Trim(re.FindAllString(string(contentModInfo), -1)[0], "\r\n")
			stringModId = strings.Replace(stringModId, "id=", "", -1)
			if !c.ContainsIgnoreModId(stringModId) {
				c.Params.ModIds = append(c.Params.ModIds, stringModId)
				fmt.Println("ModId " + stringModId + " добавлен в обработчик")
				if c.ConfigData.CopyFiles {
					go c.CopyMode(pathDirMode, c.ConfigData.ResultDirMods+"/"+stringModId, stringModId, ch)
					fmt.Println("ModId " + stringModId + " скопирован")
				}
			}
		}
	}

	count := len(c.Params.ModIds)
	for value := range ch {
		fmt.Println(value + " скопирован!")
		count -= 1
		if count == 0 {
			close(ch)
		}
	}
}

func (c *CollectorPzMods) CopyMode(source string, target string, modeName string, ch chan string) {
	c.FullCopy(source, target)
	ch <- modeName
}

func (c *CollectorPzMods) FullCopy(source string, target string) {
	sourceDir, err := os.Stat(source)
	if err != nil {
		panic(1)
	}

	if sourceDir.IsDir() {
		os.Mkdir(target, 0777)
		d, err2 := ioutil.ReadDir(source)
		if err2 != nil {
			panic("Неудалось прочитать каталог: " + source)
		}

		for _, entry := range d {
			c.FullCopy(source+"/"+entry.Name(), target+"/"+entry.Name())
		}
	} else {
		file, errFile := os.Open(source)
		if errFile != nil {
			panic("Неудалось открыть файл: " + source)
		}
		defer file.Close()

		fileOut, errFileOut := os.Create(target)
		if errFileOut != nil {
			panic("Неудалось создать файл: " + target)
		}
		defer fileOut.Close()

		_, errCopy := io.Copy(fileOut, file)
		if errCopy != nil {
			panic("Неудалось скопировать файл")
		}
	}
}

func (c *CollectorPzMods) ContainsIgnoreModId(s string) bool {
	for _, modId := range c.ConfigData.IgnoreModIds {
		if modId == s {
			return true
		}
	}
	return false
}

func (c *CollectorPzMods) GenerateFileModAndWorkshopIds() {
	infoFile, errInfoFile := os.Create("info.txt")
	if errInfoFile != nil {
		panic(errInfoFile)
	}

	var infoContent string
	infoContent = "Mods=" + strings.Join(c.Params.ModIds, ";") + "\n"
	infoContent += "Workshop=" + strings.Join(c.Params.WorkshopIds, ";") + "\n"
	infoFile.WriteString(infoContent)

	infoFile.Close()

	fmt.Println("Информация записана в файл!")

}
