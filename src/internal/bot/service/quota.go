package service

type QuotaService struct {
}

func NewQuotaService() *QuotaService {
	return &QuotaService{}
}

func (*QuotaService) GetChatQuota(int64) (int, error) {
	return 1, nil
}

func (*QuotaService) UseChatQuota(int64) error {
	return nil
}
