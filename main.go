package main

import (
	//limitorder "github.com/programcpp/okotron/limit_order"
	 "github.com/programcpp/okotron/telegram"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	//limitorder.ProcessOrders()
	telegram.Run()
}
