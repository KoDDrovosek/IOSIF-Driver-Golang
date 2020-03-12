# IOSIF golang driver 

```
go get github.com/SoGoDev/gopher-iosif
```


#### Tutorial

```go
import (
	driver "github.com/SoGoDev/IOSIF-Driver-Golang"
    "fmt"
	"log"
)


func TestHandler(key, value string) {
    fmt.Println(key, value)
}

func main () {


	opt := driver.Options{
		Topics: map[string]func(key string, value string){
			"test": TestHandler,
		},
		Url:     "http://localhost:7070",
		Periods: 30,
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

```


IOSIF https://github.com/SoGoDev/IOSIF
