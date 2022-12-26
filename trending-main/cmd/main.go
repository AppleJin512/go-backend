package main

import (
	"moonbite/trending/internal/cli"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		logrus.Fatalf("error app run: %v", err)
	}
}
