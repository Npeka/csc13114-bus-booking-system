package service

import (
	"bus-booking/payment-service/internal/model"
	"context"
)

type ConstantsService interface {
	GetBanks(ctx context.Context) ([]model.BankConstant, error)
}

type ConstantsServiceImpl struct{}

func NewConstantsService() ConstantsService {
	return &ConstantsServiceImpl{}
}

// Vietnamese banks list
var vietnameseBanks = []model.BankConstant{
	{Code: "VCB", ShortName: "Vietcombank", Name: "Ngân hàng TMCP Ngoại Thương Việt Nam"},
	{Code: "TCB", ShortName: "Techcombank", Name: "Ngân hàng TMCP Kỹ Thương Việt Nam"},
	{Code: "MB", ShortName: "MBBank", Name: "Ngân hàng TMCP Quân Đội"},
	{Code: "VTB", ShortName: "VietinBank", Name: "Ngân hàng TMCP Công Thương Việt Nam"},
	{Code: "BIDV", ShortName: "BIDV", Name: "Ngân hàng TMCP Đầu Tư và Phát Triển Việt Nam"},
	{Code: "ACB", ShortName: "ACB", Name: "Ngân hàng TMCP Á Châu"},
	{Code: "AGR", ShortName: "Agribank", Name: "Ngân hàng Nông Nghiệp và Phát Triển Nông Thôn Việt Nam"},
	{Code: "SAC", ShortName: "Sacombank", Name: "Ngân hàng TMCP Sài Gòn Thương Tín"},
	{Code: "VPB", ShortName: "VPBank", Name: "Ngân hàng TMCP Việt Nam Thịnh Vượng"},
	{Code: "TPB", ShortName: "TPBank", Name: "Ngân hàng TMCP Tiên Phong"},
	{Code: "SCB", ShortName: "SCB", Name: "Ngân hàng TMCP Sài Gòn"},
	{Code: "HDB", ShortName: "HDBank", Name: "Ngân hàng TMCP Phát Triển TP.HCM"},
	{Code: "MSB", ShortName: "MSB", Name: "Ngân hàng TMCP Hàng Hải Việt Nam"},
	{Code: "SHB", ShortName: "SHB", Name: "Ngân hàng TMCP Sài Gòn - Hà Nội"},
	{Code: "VIB", ShortName: "VIB", Name: "Ngân hàng TMCP Quốc Tế Việt Nam"},
	{Code: "OCB", ShortName: "OCB", Name: "Ngân hàng TMCP Phương Đông"},
	{Code: "EIB", ShortName: "Eximbank", Name: "Ngân hàng TMCP Xuất Nhập Khẩu Việt Nam"},
	{Code: "LPB", ShortName: "LienVietPostBank", Name: "Ngân hàng TMCP Bưu Điện Liên Việt"},
	{Code: "SEA", ShortName: "SeABank", Name: "Ngân hàng TMCP Đông Nam Á"},
	{Code: "VAB", ShortName: "VietABank", Name: "Ngân hàng TMCP Việt Á"},
	{Code: "NAB", ShortName: "NamABank", Name: "Ngân hàng TMCP Nam Á"},
	{Code: "PGB", ShortName: "PGBank", Name: "Ngân hàng TMCP Xăng Dầu Petrolimex"},
	{Code: "VCB", ShortName: "Viet Capital Bank", Name: "Ngân hàng TMCP Bản Việt"},
	{Code: "BAB", ShortName: "BacABank", Name: "Ngân hàng TMCP Bắc Á"},
	{Code: "PVC", ShortName: "PVCombank", Name: "Ngân hàng TMCP Đại Chúng Việt Nam"},
	{Code: "KLB", ShortName: "Kienlongbank", Name: "Ngân hàng TMCP Kiên Long"},
	{Code: "CAKE", ShortName: "CAKE", Name: "CAKE by VPBank"},
	{Code: "UBANK", ShortName: "Ubank", Name: "Ubank by VPBank"},
	{Code: "TIMO", ShortName: "Timo", Name: "Timo by VPBank"},
}

func (s *ConstantsServiceImpl) GetBanks(ctx context.Context) ([]model.BankConstant, error) {
	return vietnameseBanks, nil
}
