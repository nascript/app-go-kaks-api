package product

import (
	"errors"
	"kaks-cloud-web-api-task/domain/product"
	_ "kaks-cloud-web-api-task/domain/product"
	"kaks-cloud-web-api-task/infrastructure"
	"kaks-cloud-web-api-task/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

type adapter struct {
	storeService product.ProductInterface
}

// NewStoreHandler NewHandler  New handler instantiates a http handler for our store service
func NewStoreHandler(storeService product.ProductInterface) *adapter {
	return &adapter{storeService: storeService}
}

// Get product by ID
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} product.Product
// @Failure 400 {object} utils.HTTPError
// @Router /api/v1/product/{id} [get]
func (h *adapter) Get(ctx *fiber.Ctx) error {

	// Tracing
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:Get")
	defer span.End()

	id := ctx.Params("id")
	resp, err := h.storeService.Find(c, id)
	if err != nil {
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, *resp, err)
		return nil
	}
	utils.ResponseWithJSON(ctx, http.StatusOK, resp, nil)
	return nil
}

// Get all products
// @Summary Get all products
// @Description Get all products
// @Tags product
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} product.Product
// @Failure 400 {object} utils.HTTPError
// @Router /api/v1/product [get]
func (h *adapter) GetAll(ctx *fiber.Ctx) error {
	// Tracing
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:GetAll")
	defer span.End()

	filter := product.Filter{
		Page:      ctx.QueryInt("page"),
		Limit:     ctx.QueryInt("limit"),
		Latitude:  ctx.Query("latitude"),
		Longitude: ctx.Query("longitude"),
		Keyword:   ctx.Query("keyword"),
	}
	p, pagination, err := h.storeService.FindAll(c, filter)
	if err != nil {
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, []*product.Product{}, err)
		return nil
	}
	utils.ResponseWithJSON(ctx, http.StatusOK, p, nil, pagination)
	return nil
}

// Create a new product
// @Summary Create a new product
// @Description Create a new product
// @Tags product
// @Accept  json
// @Produce  json
// @Param product body product.Product true "Product"
// @Success 200 {object} product.Product
// @Failure 400 {object} utils.HTTPError
// @Router /api/v1/product [post]
func (h *adapter) Create(ctx *fiber.Ctx) error {

	// Tracing
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:Create")
	defer span.End()

	dataStore := &product.Product{}
	if err := ctx.BodyParser(&dataStore); err != nil {
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}

	// do validation
	errValidation := utils.Validate(dataStore)
	if errValidation != "" {
		slog.ErrorContext(c, "Failed to Validate api:product:Create", slog.Any("err ", errValidation))
		utils.ResponseWithJSON(ctx, http.StatusUnprocessableEntity, nil, errors.New(errValidation))
		return nil
	}

	resp, err := h.storeService.Store(c, dataStore)
	if err != nil {
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}
	utils.ResponseWithJSON(ctx, http.StatusOK, resp, nil)
	return nil

}

// Update a product by ID
// @Summary Update a product by ID
// @Description Update a product by ID
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param product body product.Product true "Product"
// @Success 200 {object} product.Product
// @Failure 400 {object} utils.HTTPError
// @Router /api/v1/product/{id} [put]
func (h adapter) Update(ctx *fiber.Ctx) error {
	// Tracing
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:Update")
	defer span.End()

	paramsID := ctx.Params("id")

	// parse data store
	dataStore := &product.Product{}
	if err := ctx.BodyParser(&dataStore); err != nil {
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}

	// do validation
	errValidation := utils.Validate(dataStore)
	if errValidation != "" {
		slog.ErrorContext(c, "Failed to Validate api:product:Create", slog.Any("err ", errValidation))
		utils.ResponseWithJSON(ctx, http.StatusUnprocessableEntity, nil, errors.New(errValidation))
		return nil
	}

	// parse id to objectID
	id, errObjectID := primitive.ObjectIDFromHex(paramsID)
	if errObjectID != nil {
		slog.ErrorContext(c, "Failed to Validate api:product:Create", slog.Any("err ", errObjectID))
		utils.ResponseWithJSON(ctx, http.StatusUnprocessableEntity, nil, errObjectID)
		return nil
	}
	// throw id
	dataStore.ID = id

	// call service
	err := h.storeService.Update(c, dataStore)
	if err != nil {
		utils.ResponseWithJSON(ctx, http.StatusBadRequest, nil, err)
		return nil
	}

	utils.ResponseWithJSON(ctx, http.StatusOK, dataStore, nil)
	return nil
}

// Delete a product by ID
// @Summary Delete a product by ID
// @Description Delete a product by ID
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} utils.HTTPError
// @Router /api/v1/product/{id} [delete]
func (h *adapter) Delete(ctx *fiber.Ctx) error {
	// Tracing
	c, span := infrastructure.Tracer().Start(ctx.UserContext(), "api:store:DeleteByID")
	defer span.End()

	id := ctx.Params("id")
	err := h.storeService.DeleteById(c, id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "Deleted successfully"})
}
