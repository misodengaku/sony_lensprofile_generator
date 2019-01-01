package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"./lensprofile"
)

func main() {
	app := cli.NewApp()
	app.Name = "LensProfileGenerator"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "Create a new profile",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Value: "LENS0000.BIN",
					Usage: "Output file name",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Lens name (ex. AF-S Micro NIKKOR 60mm f/2.8G ED, max 64byte)",
				},
				cli.IntFlag{
					Name:  "focallength, f",
					Usage: "Focal length (ex. 60mm -> 60)",
				},
				cli.Float64Flag{
					Name:  "aparture, a",
					Usage: "Maximum aparture (ex. f/2.8 -> 2.8)",
				},
				cli.IntFlag{
					Name:  "vignet_br, vg",
					Value: 0,
					Usage: "Vignetting correction brightness value",
				},
				cli.IntFlag{
					Name:  "vignet_red, vr",
					Value: 0,
					Usage: "Vignetting correction red value",
				},
				cli.IntFlag{
					Name:  "vignet_blue, vb",
					Value: 0,
					Usage: "Vignetting correction blue value",
				},
				cli.IntFlag{
					Name:  "chroma_red, cr",
					Value: 0,
					Usage: "Chromatic aberration correction red value",
				},
				cli.IntFlag{
					Name:  "chroma_blue, cb",
					Value: 0,
					Usage: "Chromatic aberration correction",
				},
				cli.IntFlag{
					Name:  "distortion, d",
					Value: 0,
					Usage: "distortion correction value",
				},
			},
			Action: func(c *cli.Context) error {
				profile := lensprofile.Profile{}
				profile.LensName = c.String("name")
				profile.FocalLength = (uint16)(c.Int("focallength"))
				profile.Aparture = c.Float64("aparture")
				profile.VignettingCorrection.Brightness = c.Int("vignet_br")
				profile.VignettingCorrection.Red = c.Int("vignet_red")
				profile.VignettingCorrection.Blue = c.Int("vignet_blue")
				profile.ChromaticAberrationCorrection.Red = c.Int("chroma_red")
				profile.ChromaticAberrationCorrection.Blue = c.Int("chroma_blue")
				profile.DistortionCorrection = c.Int("distortion")
				data, err := profile.Encode()
				if err != nil {
					return err
				}

				return ioutil.WriteFile(c.String("output"), data, 0644)
			},
		},
		{
			Name:  "print",
			Usage: "View profile",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "LENS0000.BIN",
					Usage: "Input file name",
				},
			},
			Action: func(c *cli.Context) error {
				data, err := ioutil.ReadFile(c.String("input"))
				if err != nil {
					return err
				}
				p, err := lensprofile.Unmarshal(data)
				if err != nil {
					return err
				}

				fmt.Println(c.String("input"))

				fmt.Printf("LensName:\t %s\n", p.LensName)
				fmt.Printf("FocalLength:\t %d\n", p.FocalLength)
				fmt.Printf("Aparture:\t %.2f\n", p.Aparture)
				fmt.Printf("VignettingCorrectionBrightness:\t %d\n", p.VignettingCorrection.Brightness)
				fmt.Printf("VignettingCorrectionRed:\t %d\n", p.VignettingCorrection.Red)
				fmt.Printf("VignettingCorrectionBlue:\t %d\n", p.VignettingCorrection.Blue)
				fmt.Printf("ChromaticAberrationRed:\t %d\n", p.ChromaticAberrationCorrection.Red)
				fmt.Printf("ChromaticAberrationBlue:\t %d\n", p.ChromaticAberrationCorrection.Blue)
				fmt.Printf("DistortionCorrection:\t %d\n", p.DistortionCorrection)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
