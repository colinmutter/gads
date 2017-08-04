package v201702

type DraftService struct {
	Auth
}

func NewDraftService(auth *Auth) *DraftService {
	return &DraftService{Auth: *auth}
}
