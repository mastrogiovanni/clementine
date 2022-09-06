package cartridge

import "os"

type CartridgeInfo struct {
	FileName string
	GameName string
}

type Cartridge struct {
	data []byte
	info CartridgeInfo
}

func NewCartridge() *Cartridge {
	return &Cartridge{
		data: []byte{},
		info: CartridgeInfo{},
	}
}

func (c *Cartridge) Data() []byte {
	return c.data
}

func (c *Cartridge) Load(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	c.data = data
	c.info.FileName = fileName
	c.info.GameName = string(data[0x00A0 : 0x00A0+12])
	return nil
}

func (c *Cartridge) Info() CartridgeInfo {
	return c.info
}
