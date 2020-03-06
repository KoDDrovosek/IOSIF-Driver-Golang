package IOSIF_Driver_Golang

import (
	driver "github.com/SoGoDev/IOSIF-Driver-Golang"
	"log"
)

func TestHandler(key, value string) {

}

func main() {

	opt := driver.Options{
		topics: map[string]func(key string, value string){
			"test": TestHandler,
		},
		url:     "http://localhost:7070",
		periods: 30,
	}

	d := driver.New(opt)
	err := d.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = d.Publish("test", "KEY", "VALUE")
	if err != nil {
		log.Fatal(err)
	}

}
