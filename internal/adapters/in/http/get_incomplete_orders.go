package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetOrders(c echo.Context) error {
	query := queries.NewGetIncompleteOrdersQuery()
	
	queryResponse, err := s.getIncompleteOrdersHandler.Handle(c.Request().Context(), query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
		return c.JSON(http.StatusConflict, problems.NewConflict(err.Error(), "/"))
	}

	var httpResponse = make([]servers.Order, 0, len(queryResponse.Orders))
	for _, order := range queryResponse.Orders {
		location := servers.Location{
			X: order.Location.X,
			Y: order.Location.Y,
		}

		var o = servers.Order{
			Id: order.ID,
			Location: location,
		}
		httpResponse = append(httpResponse, o)
	}

	return c.JSON(http.StatusOK, httpResponse)
}

