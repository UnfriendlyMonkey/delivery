package order

const (
	StatusEmpty     Status = ""
	StatusCreated   Status = "Created"
	StatusAssigned  Status = "Assigned"
	StatusCompleted Status = "Completed"
)

type Status string

func (s Status) Equal(target Status) bool {
	return s == target
}

func (s Status) IsEmpty() bool {
	return s == StatusEmpty
}

func (s Status) String() string {
	return string(s)
}
