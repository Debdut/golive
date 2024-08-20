package lib

import (
	"fmt"
	"os"
	"time"
)

// VERSION of Package
const VERSION = "1.2.1"

func printStartMessage(path, port string, isHTTPS string) {
	// Clear the screen.
	fmt.Print("\033[2J")
	// Move the cursor to the upper-left corner of the screen.
	fmt.Print("\033[H")
	fmt.Printf("golive\n--\n")
	fmt.Printf("Serving: %s\n", path)
	fmt.Printf("Local: http://localhost%s/\n", port)

	if isHTTPS != "" {
		fmt.Printf("Local HTTPS: https://localhost%s/\n", port)
	}
}

func printServerInformation(path, port, httpsPort string) {
	// Move to the fifth row, 1st column change if more print statements are added.
	fmt.Print("\033[5;1H")
	localIP, err := GetLocalIP()
	if err == nil && localIP != "" {
		fmt.Printf("Net: http://%s%s/\n", localIP, port)
		if httpsPort != "" {
			fmt.Printf("Net HTTPS: https://%s%s/\n", localIP, httpsPort)
		}
	} else {
		// If there is no network connection, erase the line.
		fmt.Print("\033[K")
		fmt.Println()
	}
	fmt.Println("\nRequests:", requests, "\n")
}

// Printer prints out the information associated with the server on a loop.
func Printer(dir, port, httpsPort string) {
	// Need to give time if there is a server error.
	time.Sleep(5 * time.Millisecond)
	start := time.Now()

	path, err := os.Getwd()
	// if there is an error or we are using a special path, use dir arg.
	if err != nil || dir != "./" {
		path = dir
	}

	printStartMessage(path, port, httpsPort)
	for {
		time.Sleep(500 * time.Millisecond)
		printServerInformation(path, port, httpsPort)
		// Move to the timeSince row, and clear it.
		fmt.Print("\033[8;1H")
		fmt.Print("\033[K")
		fmt.Printf("%s\n", time.Since(start).Round(time.Second))
	}
}
