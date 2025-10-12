package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateOrder(c echo.Context) error {
	createOrderCommand, err := commands.NewCreateOrderCommand(uuid.New(), "some_street", kernel.Volume(5))
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createOrderHandler.Handle(c.Request().Context(), createOrderCommand)
	if err != nil {
		c.Logger().Errorf("CreateOrder handler error: %v", err)
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
		return problems.NewConflict(err.Error(), "/")
	}

	return c.NoContent(http.StatusCreated)
}
