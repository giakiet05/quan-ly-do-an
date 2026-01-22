package util

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ProjectExcelData represents a project parsed from Excel
type ProjectExcelData struct {
	Title       string
	Description string
	Amount      int
	MinMember   int
	MaxMember   int
}

// ProjectParseError represents an error encountered while parsing a row
type ProjectParseError struct {
	Row    int
	Errors []string
}

// ParseProjectsFromExcel parses an Excel file and extracts project data
// Expected format (Vietnamese headers):
// | Tên đề tài | Mô tả | Số nhóm | Số SV tối thiểu | Số SV tối đa |
// |------------|-------|---------|-----------------|--------------|
// | Project 1  | ...   | 3       | 2               | 4            |
func ParseProjectsFromExcel(fileHeader *multipart.FileHeader) ([]ProjectExcelData, []ProjectParseError, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Excel file: %w", err)
	}
	defer f.Close()

	// Get the first sheet name
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("no sheets found in Excel file")
	}
	sheetName := sheets[0]

	// Get all rows from the first sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, nil, fmt.Errorf("Excel file must have at least 2 rows (header + data)")
	}

	// Parse projects and collect errors
	projects := make([]ProjectExcelData, 0)
	parseErrors := make([]ProjectParseError, 0)

	for i, row := range rows {
		// Skip header row (row 0)
		if i == 0 {
			continue
		}

		// Skip empty rows
		if len(row) == 0 {
			continue
		}

		// Check if row has enough columns
		if len(row) < 5 {
			parseErrors = append(parseErrors, ProjectParseError{
				Row:    i + 1,
				Errors: []string{"Thiếu dữ liệu cột (cần đủ 5 cột: Tên đề tài, Mô tả, Số nhóm, Số SV tối thiểu, Số SV tối đa)"},
			})
			continue
		}

		// Extract and trim values
		title := strings.TrimSpace(row[0])
		description := strings.TrimSpace(row[1])
		amountStr := strings.TrimSpace(row[2])
		minMemberStr := strings.TrimSpace(row[3])
		maxMemberStr := strings.TrimSpace(row[4])

		// Parse integers
		amount, errAmount := strconv.Atoi(amountStr)
		minMember, errMinMember := strconv.Atoi(minMemberStr)
		maxMember, errMaxMember := strconv.Atoi(maxMemberStr)

		// Validate and collect errors
		rowErrors := make([]string, 0)

		if title == "" {
			rowErrors = append(rowErrors, "Tên đề tài không được để trống")
		}

		if errAmount != nil || amount <= 0 {
			rowErrors = append(rowErrors, "Số nhóm phải là số nguyên dương (> 0)")
		}

		if errMinMember != nil || minMember <= 0 {
			rowErrors = append(rowErrors, "Số SV tối thiểu phải là số nguyên dương (> 0)")
		}

		if errMaxMember != nil || maxMember <= 0 {
			rowErrors = append(rowErrors, "Số SV tối đa phải là số nguyên dương (> 0)")
		}

		if errMinMember == nil && errMaxMember == nil && maxMember < minMember {
			rowErrors = append(rowErrors, "Số SV tối đa phải lớn hơn hoặc bằng Số SV tối thiểu")
		}

		// If there are errors, add to parseErrors and skip this row
		if len(rowErrors) > 0 {
			parseErrors = append(parseErrors, ProjectParseError{
				Row:    i + 1,
				Errors: rowErrors,
			})
			continue
		}

		// Add valid project to result
		projects = append(projects, ProjectExcelData{
			Title:       title,
			Description: description,
			Amount:      amount,
			MinMember:   minMember,
			MaxMember:   maxMember,
		})
	}

	if len(projects) == 0 {
		return nil, parseErrors, fmt.Errorf("no valid projects found in Excel file")
	}

	return projects, parseErrors, nil
}
