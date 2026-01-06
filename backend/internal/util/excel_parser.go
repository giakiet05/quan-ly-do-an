package util

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ParseStudentCodesFromExcel parses an Excel file and extracts student codes
// Expected format: First column contains student codes, starting from row 2 (row 1 is header)
// Example:
// | Student Code |
// |--------------|
// | 2211234      |
// | 2211235      |
func ParseStudentCodesFromExcel(fileHeader *multipart.FileHeader) ([]string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel file: %w", err)
	}
	defer f.Close()

	// Get the first sheet name
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}
	sheetName := sheets[0]

	// Get all rows from the first sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must have at least 2 rows (header + data)")
	}

	// Extract student codes from first column (skip header row)
	studentCodes := make([]string, 0)
	for i, row := range rows {
		// Skip header row (row 0)
		if i == 0 {
			continue
		}

		// Skip empty rows
		if len(row) == 0 {
			continue
		}

		// Get first column value
		studentCode := strings.TrimSpace(row[0])
		if studentCode == "" {
			continue
		}

		studentCodes = append(studentCodes, studentCode)
	}

	if len(studentCodes) == 0 {
		return nil, fmt.Errorf("no student codes found in Excel file")
	}

	return studentCodes, nil
}

// ValidateExcelFile validates if the uploaded file is a valid Excel file
func ValidateExcelFile(fileHeader *multipart.FileHeader) error {
	// Check file extension
	filename := strings.ToLower(fileHeader.Filename)
	if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") {
		return fmt.Errorf("invalid file type. Only .xlsx and .xls files are allowed")
	}

	// Check file size (max 10MB)
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file too large. Maximum size is 10MB")
	}

	return nil
}
