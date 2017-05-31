package v201601

type TargetIdeaService struct {
	Auth
}

func NewTargetIdeaService(auth *Auth) *TargetIdeaService {
	return &TargetIdeaService{Auth: *auth}
}
