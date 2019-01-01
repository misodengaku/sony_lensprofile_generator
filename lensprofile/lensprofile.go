package lensprofile

import (
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

// typedef struct {
// 	uint8_t lens_name_len; // max 64(0x40)
// 	uint8_t lens_name[64];
// 	uint8_t padding0[320]; // 0
// 	uint8_t unknown_data1; // 1, not related to focal_length field?
// 	uint8_t focal_length_msb; // 256mm -> 01, 123mm -> 00
// 	uint8_t focal_length_lsb; // 256mm -> 00, 123mm -> 7B
// 	uint8_t use_max_aparture_value; // 1 -> use max_aparture field value
// 	uint8_t max_aparture; // f7.1 -> 7
// 	uint8_t max_aparture_decimal; // f7.1 -> 10, f6.3 -> 30
// 	int8_t visnetting_brightness;
// 	int8_t visnetting_red;
// 	int8_t visnetting_blue;
// 	int8_t ca_red; // ca = Chromatic aberration
// 	int8_t ca_blue;
// 	int8_t distortion;
// 	uint8_t padding1[92]; // 0
// } LensProfile;
