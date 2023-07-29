package logger

import (
  "log"
  "fmt"
)

var okColour string = "\033[36m"
var errorColour string = "\033[31m"

func LogSimple(output string) {
  fmt.Println(output)
}

func Log(output string) {
  fileFlags()
  log.Println(okColour, output)
}

func LogV(output string) {
  log.Println(okColour, output)
}

func LogError(output string) {
  standardFlags()
  log.Println(errorColour, output)
}

func standardFlags() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func fileFlags() {
  log.SetFlags(log.Lshortfile)
}

