package config

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	File string
	CurrentPath string
}

/*
func setPath(newPath string) error {
	err := ioutil.WriteFile("fargo.txt", []byte("{ \"config\": {\"path\": newPath }}"), 0644)
	return err
}
*/

func (cfg Config) Print() {
	fmt.Println("Config:", cfg)
}

func (cfg Config) ToJson() ([]byte, error){
	b, err := json.Marshal(cfg)
	return b, err
}

func (cfg Config) Save() error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfg.File, b, 0644)
	return err
}

func Load(filename string) Config {
	cfg := Config{filename, "/"}
	
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg
	}
	
	err = json.Unmarshal(b, &cfg)
	return cfg
}
