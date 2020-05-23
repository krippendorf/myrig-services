package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var appCtx *ApplicationContext

func main() {
	var appContext ApplicationContext
	appCtx = &appContext

	readApplicationConfig()

	go startPollingRotor()

	go func() {
		router := CreateRouter()
		log.Fatal(http.ListenAndServe(appContext.ListenString, router))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}

func readApplicationConfig() {
	viper.SetConfigName("ms")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/ms/")
	viper.AddConfigPath("$HOME/.ms")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Panicf("fatal error: %s", err)
	}
	viper.Unmarshal(&appCtx)
}