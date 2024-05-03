package config

import (
	"fmt"
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	path := "config.yml"

	conf := LoadConfig(path)
	fmt.Println("config", conf)
}
