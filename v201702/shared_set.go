package v201609

type SharedSetService struct {
	Auth
}

func NewSharedSetService(auth *Auth) *SharedSetService {
	return &SharedSetService{Auth: *auth}
}
