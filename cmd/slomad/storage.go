package main

import (
	"fmt"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/sirupsen/logrus"
)

func CreateVolumes() {
	fmt.Println("Creating volumes...")
	toCreate := []string{
		"mariadb",
	}

	for _, v := range toCreate {
		fmt.Printf("Creating volume: %s\n", v)
		if err := slomad.CreateVolume(v); err != nil {
			logrus.Fatalln(err)
		}
	}
}
