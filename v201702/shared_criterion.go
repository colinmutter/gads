package v201702

type SharedCriterionService struct {
	Auth
}

func NewSharedCriterionService(auth *Auth) *SharedCriterionService {
	return &SharedCriterionService{Auth: *auth}
}
