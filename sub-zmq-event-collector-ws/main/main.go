package main

import (
	"../logger"
	"io/ioutil"
	"os"
)

func init() {
	logger.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	logger.Info.Println("test")
}
