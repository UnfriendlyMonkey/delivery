package commands

type AssignOrderCommand struct {
	isValid bool
}

func NewAssignOrderCommand() AssignOrderCommand {
	return AssignOrderCommand{
		isValid: true,
	}
}

func (c AssignOrderCommand) IsValid() bool {
	return c.isValid
}
