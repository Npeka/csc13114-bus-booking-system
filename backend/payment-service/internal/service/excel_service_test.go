package service

import (
	"testing"
	"time"

	"bus-booking/payment-service/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestNewExcelService(t *testing.T) {
	service := NewExcelService()

	assert.NotNil(t, service)
	assert.IsType(t, &ExcelServiceImpl{}, service)
}

func TestGenerateRefundExcel_Success(t *testing.T) {
	service := NewExcelService()

	refunds := []*model.RefundExportItem{
		{
			BookingReference: "BK001",
			UserName:         "Nguyễn Văn A",
			BankCode:         "VCB",
			BankName:         "Vietcombank",
			AccountNumber:    "1234567890",
			AccountHolder:    "NGUYEN VAN A",
			RefundAmount:     100000,
			Reason:           "Hủy vé do thay đổi lịch trình",
			CreatedDate:      time.Now(),
		},
		{
			BookingReference: "BK002",
			UserName:         "Trần Thị B",
			BankCode:         "TCB",
			BankName:         "Techcombank",
			AccountNumber:    "9876543210",
			AccountHolder:    "TRAN THI B",
			RefundAmount:     150000,
			Reason:           "Hủy vé do bận việc đột xuất",
			CreatedDate:      time.Now(),
		},
	}

	excelData, err := service.GenerateRefundExcel(refunds)

	assert.NoError(t, err)
	assert.NotNil(t, excelData)
	assert.Greater(t, len(excelData), 0, "Excel file should have data")

	// Verify it's a valid Excel file by checking the first few bytes (ZIP signature)
	// Excel files (.xlsx) are ZIP archives
	assert.Equal(t, byte(0x50), excelData[0], "Should start with ZIP signature 'PK'")
	assert.Equal(t, byte(0x4B), excelData[1], "Should start with ZIP signature 'PK'")
}

func TestGenerateRefundExcel_EmptyList(t *testing.T) {
	service := NewExcelService()

	refunds := []*model.RefundExportItem{}

	excelData, err := service.GenerateRefundExcel(refunds)

	assert.NoError(t, err)
	assert.NotNil(t, excelData)
	assert.Greater(t, len(excelData), 0, "Excel file should be created even for empty list")
}

func TestGenerateRefundExcel_SingleRefund(t *testing.T) {
	service := NewExcelService()

	refunds := []*model.RefundExportItem{
		{
			BookingReference: "BK001",
			UserName:         "Nguyễn Văn A",
			BankCode:         "VCB",
			BankName:         "Vietcombank",
			AccountNumber:    "1234567890",
			AccountHolder:    "NGUYEN VAN A",
			RefundAmount:     500000,
			Reason:           "Hủy vé",
			CreatedDate:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
	}

	excelData, err := service.GenerateRefundExcel(refunds)

	assert.NoError(t, err)
	assert.NotNil(t, excelData)
	assert.Greater(t, len(excelData), 0)
}

func TestGenerateRefundExcel_SpecialCharacters(t *testing.T) {
	service := NewExcelService()

	refunds := []*model.RefundExportItem{
		{
			BookingReference: "BK-VN-001",
			UserName:         "Nguyễn Thị Cẩm Tú",
			BankCode:         "MB",
			BankName:         "MBBank - Ngân hàng Quân đội",
			AccountNumber:    "0123456789",
			AccountHolder:    "NGUYEN THI CAM TU",
			RefundAmount:     1000000,
			Reason:           "Hủy vé do không thể đi do Covid-19 & các lý do khác",
			CreatedDate:      time.Now(),
		},
	}

	excelData, err := service.GenerateRefundExcel(refunds)

	assert.NoError(t, err)
	assert.NotNil(t, excelData)
	assert.Greater(t, len(excelData), 0, "Should handle Vietnamese and special characters")
}

func TestGenerateRefundExcel_MissingBankName(t *testing.T) {
	service := NewExcelService()

	refunds := []*model.RefundExportItem{
		{
			BookingReference: "BK001",
			UserName:         "Test User",
			BankCode:         "XXX",
			BankName:         "", // Empty bank name
			AccountNumber:    "1234567890",
			AccountHolder:    "TEST USER",
			RefundAmount:     100000,
			Reason:           "Test refund",
			CreatedDate:      time.Now(),
		},
	}

	excelData, err := service.GenerateRefundExcel(refunds)

	assert.NoError(t, err)
	assert.NotNil(t, excelData)
	assert.Greater(t, len(excelData), 0, "Should handle missing bank name by using bank code")
}

func TestGenerateRefundExcel_LargeDataset(t *testing.T) {
	service := NewExcelService()

	// Generate 100 refund items
	refunds := make([]*model.RefundExportItem, 100)
	for i := 0; i < 100; i++ {
		refunds[i] = &model.RefundExportItem{
			BookingReference: "BK" + string(rune('A'+i%26)) + string(rune('0'+i%10)),
			UserName:         "User " + string(rune('A'+i%26)),
			BankCode:         "VCB",
			BankName:         "Vietcombank",
			AccountNumber:    "1234567890",
			AccountHolder:    "USER NAME",
			RefundAmount:     100000 * (i + 1),
			Reason:           "Refund reason " + string(rune('A'+i%26)),
			CreatedDate:      time.Now().Add(time.Duration(-i) * time.Hour),
		}
	}

	excelData, err := service.GenerateRefundExcel(refunds)

	assert.NoError(t, err)
	assert.NotNil(t, excelData)
	assert.Greater(t, len(excelData), 0, "Should handle large datasets")
}
