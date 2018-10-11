package parameters

type DirectionFlag bool
type DeliveredFlag bool

const WithLowerTreshold DirectionFlag = true
const WithUpperTreshold DirectionFlag = false

const IncludeDelivered DeliveredFlag = false
const ExcludeDelivered DeliveredFlag = true

type GetOrdersForCourierParams struct {
	Since            int64         `form:"since"`
	Asc              DirectionFlag `form:"asc"`
	ExcludeDelivered DeliveredFlag `form:"exclude_delivered"`
}
