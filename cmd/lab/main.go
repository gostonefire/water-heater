package main

import (
	"fmt"
	"github.com/gostonefire/water-heater"
	"github.com/gostonefire/water-heater/internal/raspberry"
	"time"
)

func main() {
	device, err := raspberry.NewRaspberry("192.168.1.138", 5000, 1.0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(device *raspberry.Raspberry) { _ = device.Close() }(device)

	wh := waterheater.NewWaterHeater(device)

	temp, err := wh.GetTemp()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current temperature: %.2f\n", temp)

	var vout float64
	for i := 0.0; i <= 1.0; i += 0.1 {
		vout, err = wh.SetHeat(i)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Heat stear signal set to %.2f volt\n", vout)
		time.Sleep(3 * time.Second)
	}
}
