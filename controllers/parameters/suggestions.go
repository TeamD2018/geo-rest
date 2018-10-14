package parameters

type CouriersSuggestionParams struct {
	NamePrefix  string `json:"name_prefix"`
	PhonePrefix string `json:"phone_prefix"`
	Limit       int    `json:"limit"`
}
