# JIRA WSR Merge Tool
WSR Merge Tool is a cli tool created to verify and merge timesheets exported as
CSV files.

## Configuration
All the team and employee details should be added to config.json in the root
folder. The individual timesheet exports should be added inside `data/` folder.
The filename for each employee should as specified in the json file.

## Usage instructions
1. Export the csv data for the month's timesheet.
2. Name the csv file as shown in the `config.json` and place it inside `data/`.
3. Run `WSR.exe` and proceed.
