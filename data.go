package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	data [][]string
	dayNames = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	cellColumns = []string{"B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","AA","AB","AC","AD","AE","AF","AG"}
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func sanitizeCell(csvData [][]string) {
	for r := range csvData {
		// Remove end cell
		csvData[r] = csvData[r][:len(csvData[r]) - 1]
    if r == 0 {
      continue
    }
		// Trim Issue ID
		csvData[r][0] = strings.Split(csvData[r][0], " ")[0]
		for c := range csvData[r] {
			if c != 0 {
				hours := 0.0
				cell := csvData[r][c]
				duration := strings.Split(cell, " ")
				for _, part := range(duration) {
					if part != "" {
						if strings.HasSuffix(part, "h") {
							hour, err := strconv.Atoi(strings.TrimSuffix(part, "h"))
							if err != nil {
								log.Fatal(err)
							}
							hours += float64(hour)
						}
						if strings.HasSuffix(part, "m") {
							min, err := strconv.Atoi(strings.TrimSuffix(part, "m"))
							if err != nil {
								log.Fatal(err)
							}
							hours += float64(min)/60.0
						}
					}
				}
				csvData[r][c] = fmt.Sprint(roundFloat(hours, 2))
			}
		}
	}
}

func readFile(employeeFile string) [][]string {
	data, err := os.ReadFile("./data/" + employeeFile)
	if err != nil {
		fmt.Println("Can't read file.")
	}
	if len(data) < 100 {
		fmt.Printf("Data in %s seems to be incomplete.\n", "./data/" + employeeFile)
		return [][]string{{""}}
	}
	// Remove BOM
	data = bytes.TrimPrefix(data, []byte("\xEF\xBB\xBF"))
	csvReader := csv.NewReader(bytes.NewReader(data))
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Panic(err)
	}
	if csvData[len(csvData) - 1][0] == "_  " {
		csvData = csvData[:len(csvData) - 1]
	}
	sanitizeCell(csvData)
	return csvData
}

func cleanLeaves() {
  for _, row := range(data) {
    if row[0] == "PROJ-41" {
      for i := range row {
        if i == 0 {
          continue
        }
        if row[i] != "0" {
          row[i] = "0"
        }
      }
    }
  }
}

func addEmployeeData(employees []Employee) {
	for _, employee := range(employees) {
		data = append(data, []string{employee.Name})
		data = append(data, []string{""})
		start_pos := len(data) + 1
		csvData := readFile(employee.FileName)
		data = append (data, csvData...)
		end_pos := len(data)
		formula_row := []string{"Total"}
		for i := 0; i < len(csvData[0]); i++ {
			if i != len(csvData[0]) - 1 {
				formula_row = append(formula_row, "=SUM(" + cellColumns[i] + fmt.Sprint(start_pos) + ":"+ cellColumns[i] + fmt.Sprint(end_pos) + ")")
			} else {
				formula_row = append(formula_row, "=SUM(" + cellColumns[0] + fmt.Sprint(end_pos + 1) + ":"+ cellColumns[len(csvData[0]) - 2] + fmt.Sprint(end_pos + 1) + ")")
			}
		}
		data = append(data, formula_row)
		data = append(data, []string{""})
	}
}

func executeMerge() [][]string {
	data = append(data, []string{""})
	for _, team := range configData.Teams {
		data = append(data, []string{team.TeamName})
		data = append(data, []string{""})
		addEmployeeData(team.EmployeesList)
		data = append(data, []string{""})
	}
	return data
}
