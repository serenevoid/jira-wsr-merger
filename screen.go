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
  content := "JIRA WSR Merge Tool\n"
  content += "===================\n"
  content += "1. Generate WSR\n"
  content += "2. Verify files\n"
  content += "3. View config\n"
  content += "4. Quit\n"
  content += "-------------------\n"
  content += "Enter choice: "
  fmt.Print(content)
  fmt.Scan(&choice)
  clearInputBuffer()
  isQuitting := false
  switch choice {
  case 1:
    fmt.Println("Merging Data")
    mergedData := executeMerge()
    cleanLeaves()
    writeData(mergedData)
    isQuitting = true
    break
  case 2:
    checkIfFilesExist()
    break
  case 3:
    clearScreen()
    showConfig()
    break
  case 4:
    fmt.Println("Quitting...")
    isQuitting = true
    break
  default:
    fmt.Println("Invalid choice\nPress any key...")
    fmt.Scan()
    break
  }
  if !isQuitting {
    fmt.Println("Press any key to go back...")
    fmt.Scan()
    clearInputBuffer()
    homeScreen()
  }
}
