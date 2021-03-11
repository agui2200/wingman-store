package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var defaultConfigFile = "store.yml"

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd error: %+v", err)
	}
	defaultConfigFile = path.Join(pwd, defaultConfigFile)
}

// examples formats the given examples to the cli.
func examples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i] // indent each row with 2 spaces.
	}
	return strings.Join(ex, "\n")
}

func createDir(target string) error {
	_, err := os.Stat(target)
	if err == nil || !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return fmt.Errorf("creating schema directory: %w", err)
	}
	return nil
}
