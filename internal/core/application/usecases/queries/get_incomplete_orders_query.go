package queries

type GetIncompleteOrdersQuery struct {
	isValid bool
}

func NewGetIncompleteOrdersQuery() GetIncompleteOrdersQuery {
	return GetIncompleteOrdersQuery{
		isValid: true,
	}
}

func (q GetIncompleteOrdersQuery) IsValid() bool {
	return q.isValid
}
