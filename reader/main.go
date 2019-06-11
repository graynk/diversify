package main

import (
	"C"
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
	for _, tag := range tags {
		err = tag.Connect()
		if err != nil {
			log.Fatalf("Can't connect, %s", err)
		}
		fmt.Println(tag)
		fmt.Println(tag.UID())
		if tag.Type() != freefare.UltralightC {
			return
		}

		ultralightC := tag.(freefare.UltralightTag)
		masterKeySlice, err := hex.DecodeString(args[1])
		if err != nil || len(masterKeySlice) != 16 {
			log.Fatalf("Failed to decode hex string to 16-byte key, %s", args[1])
		}
		var masterKeyBytes [16]byte
		copy(masterKeyBytes[:], masterKeySlice)
		masterKey := freefare.NewDESFire3DESKey(masterKeyBytes)
		newKey, err := ultralightC.Diversify(*masterKey)
		if err != nil {
			log.Fatalf("Failed to diversify master key, %s", err)
		}
		err = ultralightC.Authenticate(*newKey)
		if err != nil {
			log.Fatal(err)
		}
		for page := byte(0); page < 44; page++ {
			readBytes, err := ultralightC.ReadPage(page)
			if err != nil {
				log.Fatalf("Can't read, %s", err)
			}
			fmt.Printf("% X\n", readBytes)
		}
		err = tag.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	}
}
