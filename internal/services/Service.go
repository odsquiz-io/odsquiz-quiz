package services

import (
	"errors"
	"sort"

	"github.com/kauanpecanha/odsquiz-quiz/internal/apperrors"
	"github.com/kauanpecanha/odsquiz-quiz/internal/models"
	"github.com/kauanpecanha/odsquiz-quiz/internal/repositories"
)

type Service struct {
	Repo repositories.Repository
}

func (s *Service) ListQuestions() ([]models.QuizQuestion, error) {
	return s.Repo.ListQuestions()
}

func (s *Service) Submit(userID string, request models.SubmitQuizRequest) (*models.SubmitQuizResponse, error) {
	if userID == "" {
		return nil, apperrors.Unauthorized(apperrors.CodeUserIDUnavailable, nil)
	}
	if len(request.Answers) == 0 {
		return nil, apperrors.BadRequest(apperrors.CodeIncompleteQuiz, errors.New("answers are required"))
	}

	activeQuestions, err := s.Repo.CountActiveQuestions()
	if err != nil {
		return nil, err
	}
	if int64(len(request.Answers)) != activeQuestions {
		return nil, apperrors.BadRequest(apperrors.CodeIncompleteQuiz, errors.New("all active questions must be answered"))
	}

	answers, err := s.resolveAnswers(request.Answers)
	if err != nil {
		return nil, err
	}

	odsPoints := make(map[int]int, 17)
	totalPoints := 0
	for _, answer := range answers {
		odsPoints[answer.ODSNumber] += answer.Points
		totalPoints += answer.Points
	}

	odsScores := make([]models.ODSScore, 0, len(odsPoints))
	for odsNumber, points := range odsPoints {
		odsScores = append(odsScores, models.ODSScore{ODSNumber: odsNumber, Points: points})
	}
	sort.Slice(odsScores, func(i, j int) bool {
		return odsScores[i].ODSNumber < odsScores[j].ODSNumber
	})

	submission, err := s.Repo.CreateSubmission(userID, answers, odsScores, totalPoints)
	if err != nil {
		return nil, err
	}

	return &models.SubmitQuizResponse{
		SubmissionID: submission.ID,
		UserID:       userID,
		TotalPoints:  totalPoints,
		ODSScores:    odsScores,
	}, nil
}

func (s *Service) ListSubmissions(userID string) ([]models.SubmissionSummary, error) {
	if userID == "" {
		return nil, apperrors.Unauthorized(apperrors.CodeUserIDUnavailable, nil)
	}

	submissions, err := s.Repo.ListSubmissionsByUser(userID)
	if err != nil {
		return nil, err
	}

	// Non-nil so an empty history is encoded as [] rather than null.
	summaries := make([]models.SubmissionSummary, 0, len(submissions))
	for _, submission := range submissions {
		odsScores := make([]models.ODSScore, 0, len(submission.ODSScores))
		for _, score := range submission.ODSScores {
			odsScores = append(odsScores, models.ODSScore{ODSNumber: score.ODSNumber, Points: score.Points})
		}

		summaries = append(summaries, models.SubmissionSummary{
			ID:          submission.ID,
			TotalPoints: submission.TotalPoints,
			CompletedAt: submission.CompletedAt,
			ODSScores:   odsScores,
		})
	}

	return summaries, nil
}

func (s *Service) resolveAnswers(requestAnswers []models.SubmittedAnswer) ([]models.QuizSubmissionAnswer, error) {
	hasIDs := requestAnswers[0].QuestionID != "" || requestAnswers[0].AnswerOptionID != ""
	if hasIDs {
		return s.resolveAnswersByID(requestAnswers)
	}

	return s.resolveAnswersByQuestionNumber(requestAnswers)
}

func (s *Service) resolveAnswersByID(requestAnswers []models.SubmittedAnswer) ([]models.QuizSubmissionAnswer, error) {
	selectedByQuestionID := make(map[string]string, len(requestAnswers))
	for _, answer := range requestAnswers {
		if answer.QuestionID == "" || answer.AnswerOptionID == "" {
			return nil, apperrors.BadRequest(apperrors.CodeInvalidRequest, errors.New("question_id and answer_option_id are required"))
		}
		if _, exists := selectedByQuestionID[answer.QuestionID]; exists {
			return nil, apperrors.BadRequest(apperrors.CodeDuplicateQuestion, errors.New("duplicate question"))
		}
		selectedByQuestionID[answer.QuestionID] = answer.AnswerOptionID
	}

	optionsByQuestion, questionsByID, err := s.Repo.FindOptionsByIDs(selectedByQuestionID)
	if err != nil {
		return nil, err
	}
	if len(optionsByQuestion) != len(selectedByQuestionID) {
		return nil, apperrors.BadRequest(apperrors.CodeAnswerNotFound, errors.New("invalid answer option"))
	}

	resolved := make([]models.QuizSubmissionAnswer, 0, len(requestAnswers))
	for questionID := range selectedByQuestionID {
		question, ok := questionsByID[questionID]
		if !ok {
			return nil, apperrors.BadRequest(apperrors.CodeQuestionNotFound, errors.New("invalid question"))
		}
		option := optionsByQuestion[questionID]
		resolved = append(resolved, submissionAnswer(question, option))
	}

	sort.Slice(resolved, func(i, j int) bool {
		return resolved[i].QuestionNumber < resolved[j].QuestionNumber
	})
	return resolved, nil
}

func (s *Service) resolveAnswersByQuestionNumber(requestAnswers []models.SubmittedAnswer) ([]models.QuizSubmissionAnswer, error) {
	selectedByQuestionNumber := make(map[int]int, len(requestAnswers))
	for index, answer := range requestAnswers {
		questionNumber := answer.QuestionNumber
		if questionNumber == 0 {
			questionNumber = index + 1
		}
		if answer.OptionValue < 1 || answer.OptionValue > 5 {
			return nil, apperrors.BadRequest(apperrors.CodeInvalidRequest, errors.New("option_value must be between 1 and 5"))
		}
		if _, exists := selectedByQuestionNumber[questionNumber]; exists {
			return nil, apperrors.BadRequest(apperrors.CodeDuplicateQuestion, errors.New("duplicate question"))
		}
		selectedByQuestionNumber[questionNumber] = answer.OptionValue
	}

	optionsByQuestion, questionsByNumber, err := s.Repo.FindOptionsByQuestionNumbers(selectedByQuestionNumber)
	if err != nil {
		return nil, err
	}
	if len(optionsByQuestion) != len(selectedByQuestionNumber) {
		return nil, apperrors.BadRequest(apperrors.CodeAnswerNotFound, errors.New("invalid answer option"))
	}

	resolved := make([]models.QuizSubmissionAnswer, 0, len(requestAnswers))
	for questionNumber := range selectedByQuestionNumber {
		question, ok := questionsByNumber[questionNumber]
		if !ok {
			return nil, apperrors.BadRequest(apperrors.CodeQuestionNotFound, errors.New("invalid question"))
		}
		option := optionsByQuestion[questionNumber]
		resolved = append(resolved, submissionAnswer(question, option))
	}

	sort.Slice(resolved, func(i, j int) bool {
		return resolved[i].QuestionNumber < resolved[j].QuestionNumber
	})
	return resolved, nil
}

func submissionAnswer(question models.QuizQuestion, option models.QuizAnswerOption) models.QuizSubmissionAnswer {
	return models.QuizSubmissionAnswer{
		QuestionID:     question.ID,
		AnswerOptionID: option.ID,
		QuestionNumber: question.Number,
		ODSNumber:      question.ODSNumber,
		OptionValue:    option.OptionValue,
		Points:         option.Score,
	}
}
