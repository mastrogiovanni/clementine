package main

import (
	"flag"
	"log"

	"github.com/mastrogiovanni/gameboy/cartridge"
	"github.com/mastrogiovanni/gameboy/cpu"
)

type CCart struct {
	Test uint8 `struct:"[]byte"`
	// First    []byte `struct:"[160]byte"`
	// FileName []byte `struct:"[12]byte"`
}

func main() {

	log.Println(byte(0b_1111_1111))

	log.Println("clementine v0.1.0")

	var fileName = flag.String("f", "", "Cartridge")
	flag.Parse()

	log.Println("Loading:", *fileName)

	c := cartridge.NewCartridge()
	c.Load(*fileName)
	log.Printf("Loading: '%s'\n", c.Info().GameName)

	cpu := cpu.NewCpu(c.Data())
	for {
		cpu.Step()
	}

}
