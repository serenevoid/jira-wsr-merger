package main

import (
	"fmt"
	"os"
	"os/exec"
)

var choice int

func clearInputBuffer() {
	var buf [1]byte
	_, _ = os.Stdin.Read(buf[:])
}

func clearScreen() {
	var cmd exec.Cmd
	cmd = *exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func homeScreen() {
	clearScreen()
	fmt.Println("JIRA WSR Merge Tool")
	fmt.Println("===================")
	fmt.Println("1. View config")
	fmt.Println("2. Verify files")
	fmt.Println("3. Generate WSR")
	fmt.Println("4. Quit")
	fmt.Println("-------------------")
	fmt.Print("Enter choice: ")
	fmt.Scan(&choice)
	clearInputBuffer()
	switch(choice) {
	case 1:
		clearScreen()
		showConfig()
		fmt.Println("Press any key to go back...")
		fmt.Scan()
		clearInputBuffer()
		homeScreen()
		break;
	case 2:
		clearScreen()
		fmt.Println("Verifying Files")
		checkIfFilesExist()
		fmt.Println("Press any key to go back...")
		fmt.Scan()
		clearInputBuffer()
		homeScreen()
		break;
	case 3:
		clearScreen()
		fmt.Println("Merging Data")
		mergedData := executeMerge()
		writeData(mergedData)
		break;
	case 4:
		fmt.Println("Quitting...")
		break;
	default:
		fmt.Println("Invalid choice")
		fmt.Println("Press any key...")
		fmt.Scan()
		homeScreen()
		break;
	}
}
