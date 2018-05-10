package model

import (
	"bufio"
	"config"
	"encoding/json"
	"entity"
	"os"
	"sync"
	agendaLogger "util/logger"
	log "util/logger"
)

var wg sync.WaitGroup
var fin = os.Stdin

func backupOldFiles() {
	if err := os.MkdirAll(config.BackupDir(), 0777); err != nil {
		log.Fatal(err)
	}
}

// Load : load all resources for agenda.
func Load() {
	loadConfig()
	loadAllUser()
	loadAllMeeting()
}

// Save : Save all data for agenda.
func Save() error {
	if err := saveAllUser(); err != nil {
		agendaLogger.Println(err.Error())
		return err
	}
	if err := saveAllMeeting(); err != nil {
		agendaLogger.Println(err.Error())
		return err
	}
	if err := saveConfig(); err != nil {
		agendaLogger.Println(err.Error())
		return err
	}
	return nil
}

func loadConfig() {
	fcfg, err := os.Open(config.AgendaConfigPath())
	if err != nil {
		log.Fatalf("Load config fail, for config path: %v\n", config.AgendaConfigPath())
	}

	if info, err := fcfg.Stat(); err != nil {
		log.Fatal(err)
	} else if info.Size() == 0 {
		log.Printf("File `%v` is empty.\n", fcfg.Name())
	}

	decoder := json.NewDecoder(fcfg)
	config.Load(decoder)
}
func saveConfig() error {
	fcfg, err := os.Create(config.AgendaConfigPath())
	if err != nil {
		log.Fatalf("Save config fail, for config path: %v\n", config.AgendaConfigPath())
	}

	encoder := json.NewEncoder(fcfg)
	return config.Save(encoder)
}

func loadAllUser() {
	fin, err := os.Open(config.UserDataRegisteredPath())
	if err != nil {
		log.Fatal(err)
	}

	if info, err := fin.Stat(); err != nil {
		log.Fatal(err)
	} else if info.Size() == 0 {
		log.Printf("File `%v` is empty.\n", fin.Name())
	}

	decoder := json.NewDecoder(fin)
	entity.LoadUsersAllRegistered(decoder)
}
func saveAllUser() error {
	fout, err := os.Create(config.UserDataRegisteredPath())
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(fout)
	if err := entity.SaveUsersAllRegistered(encoder); err != nil {
		log.Error(err) // TODO: hadnle ?
		return err
	}
	return nil
}

func loadAllMeeting() {
	fin, err := os.Open(config.MeetingDataPath())
	if err != nil {
		log.Fatal(err)
	}

	if info, err := fin.Stat(); err != nil {
		log.Fatal(err)
	} else if info.Size() == 0 {
		log.Printf("File `%v` is empty.\n", fin.Name())
	}

	decoder := json.NewDecoder(fin)
	entity.LoadAllMeeting(decoder)
}
func saveAllMeeting() error {
	fout, err := os.Create(config.MeetingDataPath())
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(fout)
	if err := entity.SaveAllMeeting(encoder); err != nil {
		log.Error(err) // TODO: hadnle ?
		return err
	}
	return nil
}

// .....
// ref to before

func readInput() (<-chan string, error) {
	channel := make(chan string)
	scanner := bufio.NewScanner(fin)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	wg.Add(1)
	go func() {
		for scanner.Scan() {
			channel <- scanner.Text() + "\n"
		}
		defer wg.Done()
		close(channel)
	}()

	return channel, nil
}
