package lensprofile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type Profile struct {
	LensName             string
	FocalLength          uint16
	Aparture             float64
	VignettingCorrection struct {
		Brightness int
		Red        int
		Blue       int
	}
	ChromaticAberrationCorrection struct {
		Red  int
		Blue int
	}
	DistortionCorrection int // -15~15
}

func (p Profile) Validation() error {
	if len(p.LensName) == 0 {
		return fmt.Errorf("LensName is required")
	}

	if len(p.LensName) > 64 {
		return fmt.Errorf("LensName must not be longer than 64 bytes")
	}

	if p.Aparture < 0 {
		return fmt.Errorf("Aparture must be positive value")
	}

	// TODO: ChromaticAberrationCorrection

	if p.DistortionCorrection > 15 && -15 < p.DistortionCorrection {
		return fmt.Errorf("DistortionCorrection must be in the range -15 to 15")
	}

	return nil
}

func (p Profile) Encode() ([]byte, error) {
	if err := p.Validation(); err != nil {
		return nil, fmt.Errorf("validation fail: %s", err)
	}

	data := make([]byte, 496, 496)

	// unknown value 01 00 01 00 01 ea 00
	data[0] = 0x01
	// data[1] = 0x00
	data[2] = 0x01
	// data[3] = 0x00
	data[4] = 0x01
	data[5] = 0xea
	// data[6] = 0x00

	// LensName
	data[7] = (byte)(len(p.LensName))
	copy(data[8:], []byte(p.LensName))

	// unknown value
	// data[390] = 0x01

	// FocalLength
	binary.BigEndian.PutUint16(data[393:], p.FocalLength)

	if p.Aparture != 0 {
		data[395] = 0x01
		data[396] = (byte)(math.Trunc(p.Aparture))
		data[397] = (byte)(math.Floor(((p.Aparture - math.Trunc(p.Aparture)) * 100.0) + 0.5))
	}

	data[398] = (byte)(p.VignettingCorrection.Brightness)
	data[399] = (byte)(p.VignettingCorrection.Red)
	data[400] = (byte)(p.VignettingCorrection.Blue)

	data[401] = (byte)(p.ChromaticAberrationCorrection.Red)
	data[402] = (byte)(p.ChromaticAberrationCorrection.Blue)
	data[403] = (byte)((int8)(p.DistortionCorrection))

	return data, nil
}

func Unmarshal(data []byte) (Profile, error) {
	p := Profile{}

	// LensName
	p.LensName = string(bytes.Trim(data[8:72], "\x00"))

	// FocalLength
	p.FocalLength = binary.BigEndian.Uint16(data[393:])

	if data[395] == 0x01 {
		p.Aparture = (float64(data[396]) + (float64(data[397]) / 100))
	}

	p.VignettingCorrection.Brightness = int(data[398])
	p.VignettingCorrection.Red = int(data[399])
	p.VignettingCorrection.Blue = int(data[400])

	p.ChromaticAberrationCorrection.Red = int(data[401])
	p.ChromaticAberrationCorrection.Blue = int(data[402])
	p.DistortionCorrection = int(data[403])

	if err := p.Validation(); err != nil {
		return p, fmt.Errorf("validation fail: %s", err)
	}

	return p, nil
}
