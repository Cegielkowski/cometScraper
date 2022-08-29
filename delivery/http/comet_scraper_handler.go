package http

import (
	"cometScraper/transport/request"
	"cometScraper/usecase"
	"cometScraper/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CometScraperHandler struct {
	CometScraperUC usecase.CometScraperUsecase
}

// NewCometScraperHandler will initialize the cometScrapers / resources endpoint
func NewCometScraperHandler(e *echo.Echo, cometScraperUC usecase.CometScraperUsecase) {
	handler := &CometScraperHandler{
		CometScraperUC: cometScraperUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/comet", handler.StartProcess)
	apiV1.GET("/comet/:id", handler.GetByID)
	apiV1.GET("/comet", handler.Fetch)
	apiV1.DELETE("/comet/:id", handler.Delete)
}

func (h *CometScraperHandler) StartProcess(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateCometScraperReq

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		c.Logger().Error(err)
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	uuid, err := h.CometScraperUC.StartProcess(ctx, &req)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Process Started",
		"uuid":    uuid,
	})
}

func (h *CometScraperHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	cometScraper, err := h.CometScraperUC.GetByID(ctx, id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": cometScraper})
}

func (h *CometScraperHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	cometScrapers, err := h.CometScraperUC.Fetch(ctx)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": cometScrapers})
}

func (h *CometScraperHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.CometScraperUC.Delete(ctx, id); err != nil {
		c.Logger().Error(err)
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "scraped data deleted",
	})
}
