package main

import (
  "encoding/csv"
  "fmt"
  "log"
  "os"
)

func getFileNames() []string {
  var file_list []string
  for _, team := range configData.Teams {
    for _, employee := range team.EmployeesList {
      file_list = append(file_list, employee.FileName)
    }
  }
  return file_list
}

func fileIsNotValid(filename string) bool {
  data, err := os.ReadFile("./data/" + filename)
  if err != nil {
    fmt.Println("Cannot read file. File doesn't exist or it is corrupted.")
  }
  return len(data) < 100
}

func checkIfFilesExist() {
  file_list := getFileNames()
  allFilesValid := true
  for _, filename := range(file_list) {
    if fileIsNotValid(filename) {
      fmt.Printf("%s does not exist/is not proper.\n", filename)
      allFilesValid = false
    }
  }
  if allFilesValid {
    fmt.Println("All files are valid")
  }
}

func writeData(data [][]string) {
  file, err := os.Create("./WSR_"+configData.Branch+"_"+configData.Month+".csv")
  if err != nil {
    log.Panic(err)
  }
  writer := csv.NewWriter(file)
  err = writer.WriteAll(data)
  if err != nil {
    log.Panic(err)
  }
}
