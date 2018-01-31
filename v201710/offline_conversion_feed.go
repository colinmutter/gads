package v201710

type OfflineConversionService struct {
	Auth
}

func NewOfflineConversionService(auth *Auth) *OfflineConversionService {
	return &OfflineConversionService{Auth: *auth}
}
