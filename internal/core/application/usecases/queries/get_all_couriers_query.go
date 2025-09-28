package queries

type GetAllCouriersQuery struct {
	isValid bool
}

func NewGetAllCouriersQuery() GetAllCouriersQuery {
	return GetAllCouriersQuery{
		isValid: true,
	}
}

func (q GetAllCouriersQuery) IsValid() bool {
	return q.isValid
}
