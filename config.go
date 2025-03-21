package main

import (
  "encoding/json"
  "fmt"
  "os"
)

type Employee struct {
  Name     string `json:"name"`
  FileName string `json:"filename"`
}

type Team struct {
  TeamName      string     `json:"teamname"`
  EmployeesList []Employee `json:"employees"`
}

type ConfigFile struct {
  Teams       []Team `json:"teams"`
  Holidays    []int  `json:"holidays"`
  LeaveTicket string `json:"leaveTicket"`
}

var configData ConfigFile

func loadConfig() {
  content, err := os.ReadFile("./config.json")
  if err != nil {
    fmt.Println("Cannot find/open file config.json")
  }
  json.Unmarshal(content, &configData)
}

func showConfig() {
  content := "--------------------------------------"
  for _, team := range configData.Teams {
    content += "\nTeam: " + team.TeamName + "\n"
    content += "Employees: \n"
    for _, employee := range team.EmployeesList {
      content += "  -" + employee.Name + "(file:" + employee.FileName + ")\n"
    }
    content += "--------------------------------------"
  }
  fmt.Println(content)
  fmt.Println("\nListed Holidays for the month:", configData.Holidays)
  fmt.Println("Leave Ticket ID:", configData.LeaveTicket , "\n")
}
