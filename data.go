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
  return math.Round(val * ratio) / ratio
}

func sanitizeCell(csvData [][]string) [][]string {
  i := 0;
  for i < len(csvData) {
    // Check for unused tickets
    if csvData[i][len(csvData[i]) - 1] == "" {
      csvData = append(csvData[:i], csvData[i+1:]...)
    }
    i++
  }
  for r := range csvData {
    // Remove last cell (Time Spent)
    csvData[r] = csvData[r][:len(csvData[r]) - 1]
    if r == 0 {
      continue
    }
    for c := range csvData[r] {
      if c != 0 {
        hours := 0.0
        cell := csvData[r][c]
        duration := strings.Split(cell, " ")
        for _, part := range(duration) {
          if part != "" {
            if strings.HasSuffix(part, "d") {
              hour, err := strconv.Atoi(strings.TrimSuffix(part, "d"))
              if err != nil {
                log.Fatal(err)
              }
              hours += float64(hour) * 8
            }
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
  return csvData
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
  data = bytes.TrimPrefix(data, []byte("\xEF\xBB\xBF")) // Remove BOM
  csvReader := csv.NewReader(bytes.NewReader(data))
  csvData, err := csvReader.ReadAll()
  if err != nil {
    log.Panic(err)
  }
  csvData = csvData[:len(csvData) - 1] // Remove Total row
  csvData = sanitizeCell(csvData)
  return csvData
}

func cleanLeaves() {
  leaveRow := 0
  for index, row := range(data) {
    if row[0] == configData.LeaveTicket {
      leaveRow = index
      for i := range row {
        switch row[i] {
        case "8":
          row[i] = "Full-day Leave"
          break
        case "4":
          row[i] = "Half-day Leave"
          break
        case "0":
          break
        default:
          if i != 0 {
            // If there is any other value, wrap it in
            // square brackets for better visibility
            row[i] = "[" + row[i] + "]"
          }
          break
      }
      }
      for _, day := range(configData.Holidays) {
        row[day] = "Holiday"
      }
    }
    if row[0] == "Total" && leaveRow != 0 {
      leaveRowData := data[leaveRow]
      data[leaveRow] = data[index - 1]
      data[index - 1] = leaveRowData
      leaveRow = 0;
    }
  }
}

func addEmployeeData(employees []Employee) {
  for _, employee := range(employees) {
    data = append(data, []string{employee.Name})
    data = append(data, []string{""})
    start_pos := len(data) + 2
    csvData := readFile(employee.FileName)
    data = append(data, csvData...)
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
