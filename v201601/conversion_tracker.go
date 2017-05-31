package v201601

type ConversionTrackingSettings struct {
	ConversionOptimizerMode            string `xml:conversionOptimizerMode`
	EffectiveConversionTrackingId      int64  `xml:effectiveConversionTrackingId`
	UsesCrossAccountConversionTracking bool   `xml:usesCrossAccountConversionTracking`
}

type ConversionTrackerService struct {
	Auth
}

func NewConversionTrackerService(auth *Auth) *ConversionTrackerService {
	return &ConversionTrackerService{Auth: *auth}
}
