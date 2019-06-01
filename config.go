package main

import (
	"github.com/BurntSushi/toml"
	"os"
)

func LoadWorld(path string) (*World, error) {
	world := &World{}
	if _, err := toml.DecodeFile(path, world); err != nil {
		return nil, err
	}

	return world, nil
}

func WriteWorld(path string, world *World) error {
	if file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644); err != nil {
		return err
	} else if err := toml.NewEncoder(file).Encode(world); err != nil {
		return err
	}

	return nil
}
