package main

import (
	"log"
	"fmt"
	"github.com/d2r2/go-dht"
	"flag"
	"time"
	"net/http"
	"io"
)

type Dht22 struct {
	temperature	float32
	humidity float32
	retried int
}

var (
	relayGpioPtr = flag.Int("dht22-gpio", 17, "Gpio of DHT22")
	httpPortPtr = flag.String("port", "8080", "HTTP port number")
	dht22Data Dht22
)



func gatherDht22Data(){
	for range time.Tick(time.Second * 20) {
		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithRetry(dht.DHT22, *relayGpioPtr, false, 10)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Got data")
		dht22Data = Dht22{temperature: temperature, humidity: humidity, retried: retried}
	}
}

func (dht *Dht22) get(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("%v:%v:%v", dht22Data.temperature, dht22Data.humidity, dht22Data.retried))
}

func main() {
	flag.Parse()

	go gatherDht22Data()

	mux := http.NewServeMux()
	mux.HandleFunc("/get", dht22Data.get)

	fmt.Println("Running...")
	http.ListenAndServe(":" + *httpPortPtr, mux)
}