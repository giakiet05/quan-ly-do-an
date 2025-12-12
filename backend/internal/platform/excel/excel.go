package excel

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ParseError represents an error when parsing Excel
type ParseError struct {
	Row   int    `json:"row"`
	Value string `json:"value"`
	Error string `json:"error"`
}

// ParseEmailsFromExcel parses an Excel file and extracts emails from the first column
// Expected format:
//   Row 1: Header (Email)
//   Row 2+: Email addresses
func ParseEmailsFromExcel(file *multipart.FileHeader) ([]string, []ParseError, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse Excel
	xlsx, err := excelize.OpenReader(strings.NewReader(string(fileBytes)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse Excel file: %w", err)
	}
	defer xlsx.Close()

	// Get first sheet
	sheetName := xlsx.GetSheetName(0)
	if sheetName == "" {
		return nil, nil, fmt.Errorf("no sheet found in Excel file")
	}

	// Get all rows
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("Excel file is empty")
	}

	var emails []string
	var errors []ParseError

	// Skip first row (header)
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}

		if len(row) == 0 {
			continue // Skip empty row
		}

		email := strings.TrimSpace(row[0])

		// Skip empty cells
		if email == "" {
			continue
		}

		// Validate email format (basic check)
		if !isValidEmail(email) {
			errors = append(errors, ParseError{
				Row:   i + 1,
				Value: email,
				Error: "Invalid email format",
			})
			continue
		}

		emails = append(emails, email)
	}

	return emails, errors, nil
}

// CreateInviteTemplate creates an Excel template for classroom invitations
func CreateInviteTemplate() (*excelize.File, error) {
	f := excelize.NewFile()

	// Set header
	f.SetCellValue("Sheet1", "A1", "Email")

	// Add example data
	f.SetCellValue("Sheet1", "A2", "student1@example.com")
	f.SetCellValue("Sheet1", "A3", "student2@example.com")
	f.SetCellValue("Sheet1", "A4", "student3@example.com")

	// Style header (bold, background color)
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4CAF50"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	f.SetCellStyle("Sheet1", "A1", "A1", headerStyle)

	// Set column width
	f.SetColWidth("Sheet1", "A", "A", 30)

	// Add instructions in another sheet
	f.NewSheet("Instructions")
	f.SetCellValue("Instructions", "A1", "Hướng dẫn sử dụng template")
	f.SetCellValue("Instructions", "A3", "1. Nhập email của sinh viên vào cột A (từ dòng 2 trở đi)")
	f.SetCellValue("Instructions", "A4", "2. Mỗi dòng là một email")
	f.SetCellValue("Instructions", "A5", "3. Xóa các dòng ví dụ trước khi upload")
	f.SetCellValue("Instructions", "A6", "4. Lưu file và upload lên hệ thống")

	// Set Instructions sheet column width
	f.SetColWidth("Instructions", "A", "A", 60)

	// Set Sheet1 as active
	f.SetActiveSheet(0)

	return f, nil
}

// isValidEmail does basic email validation
func isValidEmail(email string) bool {
	// Basic check: contains @ and .
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
