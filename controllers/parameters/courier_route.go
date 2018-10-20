package parameters

type CourierRoute struct {
	CourierID string `form:"courier_id"`
	Since     int64  `form:"since"`
}
