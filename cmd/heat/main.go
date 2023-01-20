package main

import (
	"fmt"
	"github.com/gostonefire/water-heater"
	"github.com/gostonefire/water-heater/internal/pid"
	"github.com/gostonefire/water-heater/internal/plot"
	"github.com/gostonefire/water-heater/internal/raspberry"
	"time"
)

func main() {
	Ts := 96.0
	pGain := 2.0
	iGain := 0.112
	dGain := 135.0

	tLog := make([][2]float64, 0)

	runTime := 300 * time.Second
	regInterval := time.Second

	// Set up system including PID regulator
	device, err := raspberry.NewRaspberry("192.168.1.138", 5000, 1.0)
	if err != nil {
		fmt.Printf("error while setting up raspberry device: %s\n", err)
		return
	}
	defer func(device *raspberry.Raspberry) { _ = device.Close() }(device)

	wh := waterheater.NewWaterHeater(device)
	reg := pid.NewPID(pGain, iGain, dGain)

	// Set up regulator interval ticker
	ticker := time.NewTicker(regInterval)
	defer ticker.Stop()
	done := make(chan bool)

	// Wait until runtime is over then signal done
	go func() {
		time.Sleep(runTime)
		done <- true
	}()

	var Th, Vd, T float64

outer:
	// Main run loop controlled by ticker
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			break outer
		case t := <-ticker.C:
			Th, err = wh.GetTemp()
			if err != nil {
				fmt.Printf("error while getting temperature from system: %s\n", err)
				return
			}
			Vd = reg.UpdatePID(Ts-Th, Th)
			_, err = wh.SetHeat(Vd)
			if err != nil {
				fmt.Printf("error while adjusting heater in system: %s\n", err)
				return
			}

			tLog = append(tLog, [2]float64{T, Th})
			fmt.Println("Current time: ", t)
		}
	}

	// Create plot from heating session
	err = plot.Plot(tLog, Ts, pGain, iGain, dGain, "C:/Develop/water-heater/tmp")
	if err != nil {
		fmt.Printf("error while creating plot: %s\n", err)
	}
}
