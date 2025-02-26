package main

import (
	"github.com/Zaman-R/otp-validator/cmd/config"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

}
