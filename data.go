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


func getMonthDetails() (int, int) {
	days := -1
	firstDay := -1
	for {
		fmt.Print("Enter number of days in the month(28-31):")
		fmt.Scan(&days)
		if days > 27 && days < 32 {
			break
		}
		fmt.Println("Invalid number of days")
	}
	for {
		fmt.Print("Enter first of day of the month(1-7):")
		fmt.Scan(&firstDay)
		if firstDay > 0 && firstDay < 8 {
			break
		}
		fmt.Println("Invalid day index")
	}
	return days, firstDay
}

func addMonthDetails(days int, day int) {
	var date_list []string
	var day_list []string
	for i := 0; i < days + 1; i++ {
		if i == 0 {
			date_list = append(date_list, "")
		} else {
			date_list = append(date_list, fmt.Sprint(i))
		}
	}
	data = append(data, date_list)
	for i := 0; i < days + 1; i++ {
		if i == 0 {
			day_list = append(day_list, "")
		} else {
			day_index := (day + i - 2) % 7
			day_list = append(day_list, dayNames[day_index])
		}
	}
	data = append(data, day_list)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func sanitizeCell(csvData [][]string) {
	for r := range csvData {
		// Remove end cell
		csvData[r] = csvData[r][:len(csvData[r]) - 1]
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
	if csvData[0][0] == "Issue" {
		csvData = csvData[1:]
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
	days, day := getMonthDetails()
	addMonthDetails(days, day)
	data = append(data, []string{""})
	for _, team := range configData.Teams {
		data = append(data, []string{team.TeamName})
		data = append(data, []string{""})
		addEmployeeData(team.EmployeesList)
		data = append(data, []string{""})
	}
	return data
}
