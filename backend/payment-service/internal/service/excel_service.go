package service

import (
	"bus-booking/payment-service/internal/model"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type ExcelService interface {
	GenerateRefundExcel(refunds []*model.RefundExportItem) ([]byte, error)
}

type ExcelServiceImpl struct{}

func NewExcelService() ExcelService {
	return &ExcelServiceImpl{}
}

func (s *ExcelServiceImpl) GenerateRefundExcel(refunds []*model.RefundExportItem) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()

	sheetName := "Danh Sách Hoàn Tiền"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}

	// Set as active sheet
	f.SetActiveSheet(index)
	// Delete default sheet
	if err = f.DeleteSheet("Sheet1"); err != nil {
		log.Error().Err(err).Msg("failed to delete default sheet")
	}

	// Define header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	// Define data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data style: %w", err)
	}

	// Define currency style for amount column
	currencyStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		CustomNumFmt: strPtr("#,##0"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create currency style: %w", err)
	}

	// Set column widths
	columnWidths := map[string]float64{
		"A": 8,  // STT
		"B": 18, // Mã đặt vé
		"C": 25, // Tên khách hàng
		"D": 20, // Ngân hàng
		"E": 20, // Số tài khoản
		"F": 25, // Chủ tài khoản
		"G": 18, // Số tiền (VND)
		"H": 35, // Lý do
		"I": 18, // Ngày tạo
	}
	for col, width := range columnWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, fmt.Errorf("failed to set column width: %w", err)
		}
	}

	// Set header
	headers := []string{
		"STT",
		"Mã Đặt Vé",
		"Tên Khách Hàng",
		"Ngân Hàng",
		"Số Tài Khoản",
		"Chủ Tài Khoản",
		"Số Tiền (VND)",
		"Lý Do Hoàn Tiền",
		"Ngày Tạo",
	}

	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, fmt.Errorf("failed to set header: %w", err)
		}
		if err := f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return nil, fmt.Errorf("failed to set header style: %w", err)
		}
	}

	// Set row height for header
	if err := f.SetRowHeight(sheetName, 1, 25); err != nil {
		return nil, fmt.Errorf("failed to set row height: %w", err)
	}

	// Fill data
	for i, refund := range refunds {
		row := i + 2
		rowStr := strconv.Itoa(row)

		// STT
		if err := f.SetCellValue(sheetName, "A"+rowStr, i+1); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Mã đặt vé
		if err := f.SetCellValue(sheetName, "B"+rowStr, refund.BookingReference); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Tên khách hàng
		if err := f.SetCellValue(sheetName, "C"+rowStr, refund.UserName); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Ngân hàng
		bankDisplay := refund.BankName
		if bankDisplay == "" {
			bankDisplay = refund.BankCode
		}
		if err := f.SetCellValue(sheetName, "D"+rowStr, bankDisplay); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Số tài khoản
		if err := f.SetCellValue(sheetName, "E"+rowStr, refund.AccountNumber); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Chủ tài khoản
		if err := f.SetCellValue(sheetName, "F"+rowStr, refund.AccountHolder); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Số tiền
		if err := f.SetCellValue(sheetName, "G"+rowStr, refund.RefundAmount); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}
		if err := f.SetCellStyle(sheetName, "G"+rowStr, "G"+rowStr, currencyStyle); err != nil {
			return nil, fmt.Errorf("failed to set cell style: %w", err)
		}

		// Lý do
		if err := f.SetCellValue(sheetName, "H"+rowStr, refund.Reason); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Ngày tạo
		dateStr := refund.CreatedDate.Format("02/01/2006 15:04")
		if err := f.SetCellValue(sheetName, "I"+rowStr, dateStr); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}

		// Apply data style to all cells in this row
		for col := 'A'; col <= 'I'; col++ {
			cell := string(col) + rowStr
			if col != 'G' { // G already has currency style
				if err := f.SetCellStyle(sheetName, cell, cell, dataStyle); err != nil {
					return nil, fmt.Errorf("failed to set cell style: %w", err)
				}
			}
		}
	}

	// Add summary row
	summaryRow := len(refunds) + 3
	summaryRowStr := strconv.Itoa(summaryRow)

	// Total label
	if err := f.SetCellValue(sheetName, "F"+summaryRowStr, "TỔNG CỘNG:"); err != nil {
		return nil, fmt.Errorf("failed to set summary label: %w", err)
	}

	// Calculate total
	totalAmount := 0
	for _, refund := range refunds {
		totalAmount += refund.RefundAmount
	}

	if err := f.SetCellValue(sheetName, "G"+summaryRowStr, totalAmount); err != nil {
		return nil, fmt.Errorf("failed to set total amount: %w", err)
	}

	// Style summary row
	summaryStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create summary style: %w", err)
	}

	if err := f.SetCellStyle(sheetName, "F"+summaryRowStr, "F"+summaryRowStr, summaryStyle); err != nil {
		return nil, fmt.Errorf("failed to set summary style: %w", err)
	}

	totalStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "FF0000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		CustomNumFmt: strPtr("#,##0"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create total style: %w", err)
	}

	if err := f.SetCellStyle(sheetName, "G"+summaryRowStr, "G"+summaryRowStr, totalStyle); err != nil {
		return nil, fmt.Errorf("failed to set total style: %w", err)
	}

	// Generate file in memory
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

func strPtr(s string) *string {
	return &s
}
