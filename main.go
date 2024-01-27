package main

import (
	"log"
	"time"

	"github.com/fishnix/environ/internal/ds18b20"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	if err := r.Connect(); err != nil {
		log.Fatalln("failed to connect adapter", err)
	}

	defer func() {
		if err := r.Finalize(); err != nil {
			log.Fatalln("failed to finalize adaptor:", err)
		}
	}()

	ds18bs20 := ds18b20.NewThermalProbeDriver("probe")
	defer func() {
		if err := ds18bs20.Halt(); err != nil {
			log.Fatalln("failed to halt probe:", err)
		}
	}()

	work := func() {
		gobot.Every(3*time.Second, func() {
			t, err := ds18bs20.ReadTempC()
			if err != nil {
				log.Println("failed to read temp:", err)
			}

			log.Printf("temp: %fC", t)
		})
	}

	robot := gobot.NewRobot("tempBot",
		[]gobot.Connection{r},
		[]gobot.Device{ds18bs20},
		work,
	)

	if err := robot.Start(); err != nil {
		log.Fatalln("failed to start robot:", err)
	}
}
