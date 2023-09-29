package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Employee struct {
	Name string `json:"name"`
	FileName string `json:"filename"`
}

type Team struct {
	TeamName string `json:"team"`
	EmployeesList []Employee `json:"employees"`
}

type ConfigFile struct {
	Conf []Team `json:"conf"`
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
	fmt.Println("--------------------------------------")
	for _, team := range configData.Conf {
		fmt.Println("Team: ", team.TeamName)
		fmt.Println("Employees: ")
		for _, employee := range team.EmployeesList {
			fmt.Println("  -", employee.Name, "(file:", employee.FileName, ")")
		}
		fmt.Println("--------------------------------------")
	}
}
