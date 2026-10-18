package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/kauanpecanha/odsquiz-quiz/internal/apperrors"
	"github.com/kauanpecanha/odsquiz-quiz/internal/models"
	"github.com/kauanpecanha/odsquiz-quiz/internal/services"
)

type Handler struct {
	Service *services.Service
}

func (h *Handler) ListQuestions(c fiber.Ctx) error {
	questions, err := h.Service.ListQuestions()
	if err != nil {
		return respondError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(questions)
}

func (h *Handler) Submit(c fiber.Ctx) error {
	request, err := parseSubmitRequest(c.Body())
	if err != nil {
		return respondError(c, apperrors.BadRequest(apperrors.CodeInvalidRequest, err))
	}

	userID, _ := c.Locals("userID").(string)
	response, err := h.Service.Submit(userID, request)
	if err != nil {
		return respondError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *Handler) ListSubmissions(c fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	submissions, err := h.Service.ListSubmissions(userID)
	if err != nil {
		return respondError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(submissions)
}

func parseSubmitRequest(body []byte) (models.SubmitQuizRequest, error) {
	var request models.SubmitQuizRequest
	if err := json.Unmarshal(body, &request); err == nil && request.Answers != nil {
		return request, nil
	}

	var compact struct {
		Answers []int `json:"answers"`
	}
	if err := json.Unmarshal(body, &compact); err != nil {
		return models.SubmitQuizRequest{}, err
	}

	answers := make([]models.SubmittedAnswer, 0, len(compact.Answers))
	for index, optionValue := range compact.Answers {
		answers = append(answers, models.SubmittedAnswer{
			QuestionNumber: index + 1,
			OptionValue:    optionValue,
		})
	}

	return models.SubmitQuizRequest{Answers: answers}, nil
}

func respondError(c fiber.Ctx, err error) error {
	appErr, ok := apperrors.From(err)
	if !ok {
		internalErr, _ := apperrors.From(apperrors.Internal(err))
		appErr = internalErr
	}

	status := fiber.StatusInternalServerError
	switch appErr.Kind {
	case apperrors.KindBadRequest:
		status = fiber.StatusBadRequest
	case apperrors.KindUnauthorized:
		status = fiber.StatusUnauthorized
	case apperrors.KindNotFound:
		status = fiber.StatusNotFound
	}

	return c.Status(status).JSON(fiber.Map{
		"code": appErr.Code,
	})
}
