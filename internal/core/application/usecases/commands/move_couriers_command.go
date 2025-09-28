package commands

type MoveCouriersCommand struct {
	isValid bool
}

func NewMoveCouriersCommand() MoveCouriersCommand {
	return MoveCouriersCommand{
		isValid: true,
	}
}

func (c MoveCouriersCommand) IsValid() bool {
	return c.isValid
}
