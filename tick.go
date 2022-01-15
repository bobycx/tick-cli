package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/inancgumus/screen"
	"github.com/mattn/go-tty"
)

var screenWidth, screenHeight = screen.Size()
var pause = false
var asciiNums = [11]string{

	`
000000
00  00
00  00
00  00
000000
	`,

	`
1111  
  11  
  11  
  11  
111111
	`,

	`
222222
     2
222222
2     
222222
	`,

	`
333333
    33
333333
    33
333333
	`,

	`
44  44
44  44
444444
    44
    44
	`,

	`
555555
55    
555555
    55
555555
	`,

	`
666666
66    
666666
66  66
666666
	`,

	`
777777
    77
    77
    77
    77
	`,

	`
888888
88  88
888888
88  88
888888
	`,

	`
999999
99  99
999999
    99
999999
	`,

	`
  
::
  
::
  
	`,
}

func digit(num, place int) int {
	r := num % int(math.Pow(10, float64(place)))
	return r / int(math.Pow(10, float64(place-1)))
}

func formatter(h int, m int, s int) []string {

	var strHours = strconv.Itoa(h)
	var strMinutes = strconv.Itoa(m)
	var strSeconds = strconv.Itoa(s)

	var asciiStrings []string

	if len(strHours) == 1 {
		asciiStrings = append(asciiStrings, asciiNums[0])
		asciiStrings = append(asciiStrings, asciiNums[h])
	} else {
		for i := 0; i < len(strHours); i++ {
			asciiStrings = append(asciiStrings, asciiNums[digit(h, len(strHours)-i)])
		}
	}

	asciiStrings = append(asciiStrings, asciiNums[10])

	if len(strMinutes) == 1 {
		asciiStrings = append(asciiStrings, asciiNums[0])
		asciiStrings = append(asciiStrings, asciiNums[m])
	} else {
		asciiStrings = append(asciiStrings, asciiNums[digit(m, 2)])
		asciiStrings = append(asciiStrings, asciiNums[digit(m, 1)])
	}

	asciiStrings = append(asciiStrings, asciiNums[10])

	if len(strSeconds) == 1 {
		asciiStrings = append(asciiStrings, asciiNums[0])
		asciiStrings = append(asciiStrings, asciiNums[s])
	} else {
		asciiStrings = append(asciiStrings, asciiNums[digit(s, 2)])
		asciiStrings = append(asciiStrings, asciiNums[digit(s, 1)])
	}

	return asciiStrings
}

func asciiConcat(strArray []string, separator string) string {

	var linesArray [][]string
	var finalResult []string

	// calculate vertical and horizontal padding needed to center display
	w, h := screen.Size()
	if w != screenWidth {
		screenWidth = (w - 54) / 2
		screen.Clear()
	} else {
		screenWidth = (screenWidth - 54) / 2
	}

	if h != screenHeight {
		screenHeight = (h - 11) / 2
		screen.Clear()
	} else {
		screenHeight = (screenHeight - 11) / 2
	}

	// append vertical padding
	finalResult = append(finalResult, strings.Repeat("\n", screenHeight))

	for i := 0; i < len(strArray); i++ {
		linesArray = append(linesArray, strings.Split(strArray[i], "\n"))
	}

	for i := 0; i < 7; i++ {
		// add horizontal padding
		var strLine = strings.Repeat(" ", screenWidth)
		for m := 0; m < len(linesArray); m++ {
			strLine += linesArray[m][i]
			strLine += separator
		}

		finalResult = append(finalResult, strLine)

	}

	return strings.Join(finalResult, "\n")

}

func detectPause() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if r == 'p' {
			if pause == true {
				pause = false
			} else {
				pause = true
			}

		}
	}
}

func catchSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()
}

func main() {

	var hours = 00
	var minutes = 00
	var seconds = 00

	// hides the console cursor
	fmt.Print("\033[?25l")

	screen.Clear()
	go detectPause()
	catchSignal()

	for {
		for pause != true {

			if seconds == 60 {
				seconds = 0
				minutes += 1
			}

			if minutes == 60 {
				minutes = 0
				hours += 1
			}
			screen.MoveTopLeft()
			fmt.Printf("\r%s", asciiConcat(formatter(hours, minutes, seconds), "  "))
			seconds += 1

			time.Sleep(time.Second)
		}
	}
}
