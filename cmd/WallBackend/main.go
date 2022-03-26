package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/katelinlis/Wallbackend/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "/configs/apiserver-prod.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()

	path, err := os.Getwd()
	fmt.Println(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = toml.DecodeFile(path+configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	//err = sentry.Init(sentry.ClientOptions{
	//	Dsn: "https://2d845747eaa34c2fa80df5a2f7c17760@o1036815.ingest.sentry.io/6004432",
	//})
	/*if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}*/

	// Flush buffered events before the program terminates.
	//defer sentry.Flush(2 * time.Second)

	//sentry.CaptureMessage("It works!")

	go apiserver.Start(config)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
}
