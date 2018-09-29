package parameters

type GetOrdersForCourierParams struct {
	Since int64 `form:"since"`
	Asc   bool  `form:"asc"`
}
