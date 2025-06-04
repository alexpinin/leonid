package service

type QuotaService struct {
}

func NewQuotaService() *QuotaService {
	return &QuotaService{}
}

func (*QuotaService) UseChatQuota(int64) bool {
	return true
}
