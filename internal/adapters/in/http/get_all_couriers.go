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

func (s *Server) GetCouriers(c echo.Context) error {
	query := queries.NewGetAllCouriersQuery()

	resp, err := s.getAllCouriersHandler.Handle(c.Request().Context(), query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
		return c.JSON(http.StatusConflict, problems.NewConflict(err.Error(), "/"))
	}

	var httpResponse = make([]servers.Courier, 0, len(resp.Couriers))
	for _, courier := range resp.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var c = servers.Courier{
			Id: courier.ID,
			Name: courier.Name,
			Location: location,
		}
		httpResponse = append(httpResponse, c)
	}

	return c.JSON(http.StatusOK, httpResponse)
}
