package main

import (
	"fmt"
	"io/ioutil"

	"./lensprofile"
)

func main() {
	profile := lensprofile.Profile{}
	profile.LensName = "AF-S Micro NIKKOR 60mm f/2.8G ED"
	profile.FocalLength = 60
	profile.Aparture = 2.8
	data, err := profile.Encode()
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}
	fmt.Printf("%#v", data)

	ioutil.WriteFile("testlens.bin", data, 0644)
}
