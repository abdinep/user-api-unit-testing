package main

import (
	"project/initialzer"
	"project/router"
)

func init() {
	initialzer.InitSetup()
}
func main() {
	server := router.Ginsetup()
	server.Run(":8080")
}
