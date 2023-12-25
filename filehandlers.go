package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/udhos/equalfile"
)

func bkpPlayerMp3() {
	paths := getVTOLDir()

	if paths == nil {
		log.Panic(fmt.Errorf("no VTOL VR Directories found. Is game installed?"))
	}

	for _, x := range paths {
		files, err := os.ReadDir(x + "RadioMusic\\")
		if err != nil {
			log.Panic(err)
		}

		for _, y := range files {

			if invalidfile := compareMp3(y, x); invalidfile {
				continue
			}
			backup(y, x)

		}

	}
}

func backup(file fs.DirEntry, path string) {

	if validBackup(file.Name()) {

		fmt.Println("File: " + file.Name() + " valid for backup.")
		os.Rename(path+"RadioMusic\\"+file.Name(), path+"RadioMusic\\"+file.Name()+".bkp")

	} else {

		fmt.Println("File: " + file.Name() + " previously backed up, not renaming.")

	}

}

func validBackup(name string) bool {
	return !strings.Contains(name, ".bkp") && !(strings.Contains(name, "0.mp3") || strings.Contains(name, "1.mp3") || strings.Contains(name, "2.mp3"))
}

func compareMp3(file fs.DirEntry, path string) bool {
	switch file.Name() {
	case "0.mp3":
		return compareHandler(path, file.Name())
	case "1.mp3":
		return compareHandler(path, file.Name())
	case "2.mp3":
		return compareHandler(path, file.Name())
	default:
		return false
	}
}

func compareHandler(path string, name string) bool {
	blank, err := embeds.Open("blank.mp3")
	if err != nil {
		log.Panic(err)
	}
	defer blank.Close()

	compare := equalfile.New(nil, equalfile.Options{})

	file, err := os.Open(path + "RadioMusic\\" + name)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	equal, err := compare.CompareReader(blank, file)
	if err != nil {
		log.Panic(err)
	}

	if !equal {
		fmt.Println("Found " + name + " in directory, detected file not equal, backing up.")
		file.Close()
		err := os.Rename(path+"RadioMusic\\"+name, path+"RadioMusic\\"+name+".bkp")
		if err != nil {
			log.Panic(err)
		}
		return false
	}
	fmt.Println("Found " + name + " in directory, detected file equals, not backing up.")
	return true
}

func InitMP3() []string {

	blank, err := embeds.ReadFile("blank.mp3")
	if err != nil {
		log.Panic(err)
	}

	paths := getVTOLDir()

	for _, x := range paths {
		err = os.WriteFile(x+"RadioMusic\\0.mp3", blank, os.ModeAppend)
		if err != nil {
			log.Panic(err)
		}
		err = os.WriteFile(x+"RadioMusic\\1.mp3", blank, os.ModeAppend)
		if err != nil {
			log.Panic(err)
		}
		err = os.WriteFile(x+"RadioMusic\\2.mp3", blank, os.ModeAppend)
		if err != nil {
			log.Panic(err)
		}
	}

	os.WriteFile("", nil, os.ModeAppend)

	return paths

}
