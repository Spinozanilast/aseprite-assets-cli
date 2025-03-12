package utils

import (
	"fmt"
	"strings"
	"time"
)

func CreateSpinner(charsSet string, stop chan bool, message string) {
	fmt.Printf("\r%s %s", message, " ")
	i := 0
	for {
		select {
		case <-stop:
			fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+2))
			stop <- true
			return
		default:
			fmt.Printf("\r%s %s%s", message, string([]rune(charsSet)[i]), strings.Repeat(" ", 5))
			i = (i + 1) % len([]rune(charsSet))
			time.Sleep(100 * time.Millisecond)
		}
	}
}
