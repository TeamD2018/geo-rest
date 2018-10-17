package parameters

type Suggestion struct {
	Prefix string `form:"prefix"`
	Limit  int    `form:"limit"`
}
