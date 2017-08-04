package v201609

type ConversionTrackingSettings struct {
	EffectiveConversionTrackingId      int64 `xml:effectiveConversionTrackingId`
	UsesCrossAccountConversionTracking bool  `xml:usesCrossAccountConversionTracking`
}

type ConversionTrackerService struct {
	Auth
}

func NewConversionTrackerService(auth *Auth) *ConversionTrackerService {
	return &ConversionTrackerService{Auth: *auth}
}
