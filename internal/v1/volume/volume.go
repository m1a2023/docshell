package volume

import (
	doconf "docshell/internal/v1/config"
	"log"
	"os"
)

// Volume path
var path string

// Set up
func init() {
	cfg := doconf.Config
	path = setupVolume(cfg.Volume)
}

// Setup volume for saving documents
func setupVolume(path string) string {
	// Create volume
	log.Printf("out %s", path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		msg := "Unable to create volume %s, because of %v"
		log.Fatalf(msg, path, err)
	}
	return path
}

func GetPath() string {
	return path
}
