package main

import (
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

func (a *App) sacnWorker() {
	ticker := time.NewTicker(time.Second)
	if a.sacnConfig.Fps != 0 {
		ticker = time.NewTicker(time.Second / time.Duration(a.sacnConfig.Fps))
	}

	work := func() {
		a.ensureSACNSender()
		a.ensureSACNUniverses()

		for uni := range a.activeUniverses {
			p := packet.NewDataPacket()
			data := a.universeDMXData[uni]
			p.SetData(data[:])
			a.activeUniverses[uni] <- p
		}
	}

	for {
		select {
		case <-a.sacnUpdatedConfig:
			ticker = time.NewTicker(time.Second)
			if a.sacnConfig.Fps != 0 {
				ticker = time.NewTicker(time.Second / time.Duration(a.sacnConfig.Fps))
			}
		case <-a.sacnStopLoop:
			a.closeSACNSender()
			a.sacnWorkerWG.Done()
			return
		case <-ticker.C:
			work()
		}
	}
}

func (a *App) ensureSACNSender() error {
	if a.sender != nil {
		return nil
	}

	sourceName := "Folje"
	hostname, err := os.Hostname()
	if err == nil {
		sourceName += "-" + hostname
	}

	opts := sacn.SenderOptions{
		SourceName: sourceName,
	}
	sender, err := sacn.NewSender(a.sacnConfig.IpAddress, &opts)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		a.sender = nil
		a.AlertDialog("sACN sender error", err.Error())
		return err
	}

	a.sender = sender

	return nil
}

func (a *App) closeSACNSender() {
	if a.sender == nil {
		return
	}

	a.sender.Close()
	a.sender = nil

	a.activeUniverses = make(map[uint16]chan<- packet.SACNPacket)
}

func (a *App) ensureSACNUniverses() error {
	err := a.ensureSACNSender()
	if err != nil {
		return err
	}

	for uni := range a.activeUniverses {
		inUse := false
		for _, fixture := range a.fixtures {
			if uni == fixture.Universe {
				inUse = true
				break
			}
		}

		if inUse {
			continue
		}

		a.deactiveUniverse(uni)
	}

	for _, fixture := range a.fixtures {
		a.activateUniverse(fixture.Universe)
	}

	return nil
}

func (a *App) activateUniverse(uni uint16) {
	if a.activeUniverses[uni] != nil {
		return
	}

	universe, err := a.sender.StartUniverse(uni)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}
	a.sender.SetMulticast(uni, a.sacnConfig.Multicast)
	for _, dest := range a.sacnConfig.Destinations {
		a.sender.AddDestination(uni, dest)
	}
	a.activeUniverses[uni] = universe
}

func (a *App) deactiveUniverse(uni uint16) {
	if a.activeUniverses[uni] == nil {
		return
	}

	a.sender.StopUniverse(uni)
	delete(a.activeUniverses, uni)
}

func (a *App) SetSACNConfig(sacnConfig SACNConfig) {
	a.sacnConfig = &sacnConfig

	if a.sender != nil {
		a.closeSACNSender()
	}

	a.ensureSACNSender()
	a.ensureSACNUniverses()

	a.sacnUpdatedConfig <- true
}

func (a *App) GetSACNConfig() SACNConfig {
	return *a.sacnConfig
}
