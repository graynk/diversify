package main

import (
	"encoding/hex"
	"fmt"
	"github.com/fuzxxl/nfc/2.0/nfc"
	"github.com/graynk/freefare/0.3-diversify/freefare"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatal("No master key provided")
	}
	devices, err := nfc.ListDevices()
	if err != nil || len(devices) == 0 {
		log.Fatalf("No devices found, %s", err)
	}
	device, err := nfc.Open(devices[0])
	if err != nil {
		log.Fatalf("Failed to open device, %s", err)
	}
	tags, err := freefare.GetTags(device)
	if err != nil {
		log.Fatalf("Failed to get tags, %s", err)
	}
	masterKeySlice, err := hex.DecodeString(args[1])
	if err != nil || len(masterKeySlice) != 16 {
		log.Fatalf("Failed to decode hex string to 16-byte key, %s", args[1])
	}
	var masterKeyBytes [16]byte
	copy(masterKeyBytes[:], masterKeySlice)
	masterKey := freefare.NewDESFire3DESKey(masterKeyBytes)

	var defaultKeyBytes [16]byte
	if len(args) > 2 {
		defaultKeySlice, err := hex.DecodeString(args[2])
		if err != nil || len(masterKeySlice) != 16 {
			log.Fatalf("Failed to decode hex string to 16-byte key, %s", args[2])
		}
		copy(defaultKeyBytes[:], defaultKeySlice)
	} else {
		// BREAKMEIFYOUCAN!
		defaultKeyBytes = [16]byte{0x49, 0x45, 0x4D, 0x4B, 0x41, 0x45, 0x52, 0x42, 0x21, 0x4E, 0x41, 0x43, 0x55, 0x4F, 0x59, 0x46}
	}
	defaultKey := freefare.NewDESFire3DESKey(defaultKeyBytes)
	for _, tag := range tags {
		err = tag.Connect()
		if err != nil {
			log.Fatalf("Failed to connect, %s", err)
		}
		fmt.Println(tag)
		fmt.Println(tag.UID())
		if tag.Type() != freefare.UltralightC {
			log.Fatal("Not Ultralight C")
		}
		ultralightC := tag.(freefare.UltralightTag)
		diversifiedKey, err := ultralightC.Diversify(*masterKey)
		if err != nil {
			log.Fatalf("Failed to diversify master key, %s", err)
		}
		err = ultralightC.SwapKeys(*defaultKey, *diversifiedKey)
		if err != nil {
			log.Fatalf("Failed to swap, %s", err)
		}
		err = ultralightC.Authenticate(*diversifiedKey)
		if err != nil {
			log.Fatalf("Failed to auth with diversified key, %s", err)
		}
		err = tag.Disconnect()
		if err != nil {
			log.Fatalf("Failed to disconnect, %s", err)
		}
		fmt.Println("Done")
	}
}
