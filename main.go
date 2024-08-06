package main

import (
	"github.com/programcpp/okotron/telegram"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	telegram.Run()
}
