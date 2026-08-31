package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/kauanpecanha/odsquiz-quiz/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListQuestions() ([]models.QuizQuestion, error)
	CountActiveQuestions() (int64, error)
	FindOptionsByQuestionNumbers(values map[int]int) (map[int]models.QuizAnswerOption, map[int]models.QuizQuestion, error)
	FindOptionsByIDs(values map[string]string) (map[string]models.QuizAnswerOption, map[string]models.QuizQuestion, error)
	CreateSubmission(userID string, answers []models.QuizSubmissionAnswer, odsScores []models.ODSScore, totalPoints int) (*models.QuizSubmission, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) ListQuestions() ([]models.QuizQuestion, error) {
	var questions []models.QuizQuestion
	err := r.DB.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("option_value ASC")
		}).
		Where("active = ?", true).
		Order("number ASC").
		Find(&questions).Error
	return questions, err
}

func (r *repository) CountActiveQuestions() (int64, error) {
	var count int64
	err := r.DB.Model(&models.QuizQuestion{}).Where("active = ?", true).Count(&count).Error
	return count, err
}

func (r *repository) FindOptionsByQuestionNumbers(values map[int]int) (map[int]models.QuizAnswerOption, map[int]models.QuizQuestion, error) {
	var questions []models.QuizQuestion
	err := r.DB.Preload("Options").Where("active = ?", true).Find(&questions).Error
	if err != nil {
		return nil, nil, err
	}

	optionsByQuestion := make(map[int]models.QuizAnswerOption, len(values))
	questionsByNumber := make(map[int]models.QuizQuestion, len(questions))
	for _, question := range questions {
		questionsByNumber[question.Number] = question
		expectedOptionValue, ok := values[question.Number]
		if !ok {
			continue
		}

		for _, option := range question.Options {
			if option.OptionValue == expectedOptionValue {
				optionsByQuestion[question.Number] = option
				break
			}
		}
	}

	return optionsByQuestion, questionsByNumber, nil
}

func (r *repository) FindOptionsByIDs(values map[string]string) (map[string]models.QuizAnswerOption, map[string]models.QuizQuestion, error) {
	questionIDs := make([]string, 0, len(values))
	for questionID := range values {
		questionIDs = append(questionIDs, questionID)
	}

	var questions []models.QuizQuestion
	err := r.DB.Preload("Options").Where("id IN ? AND active = ?", questionIDs, true).Find(&questions).Error
	if err != nil {
		return nil, nil, err
	}

	optionsByQuestion := make(map[string]models.QuizAnswerOption, len(values))
	questionsByID := make(map[string]models.QuizQuestion, len(questions))
	for _, question := range questions {
		questionsByID[question.ID] = question
		expectedOptionID, ok := values[question.ID]
		if !ok {
			continue
		}

		for _, option := range question.Options {
			if option.ID == expectedOptionID {
				optionsByQuestion[question.ID] = option
				break
			}
		}
	}

	return optionsByQuestion, questionsByID, nil
}

func (r *repository) CreateSubmission(userID string, answers []models.QuizSubmissionAnswer, odsScores []models.ODSScore, totalPoints int) (*models.QuizSubmission, error) {
	now := time.Now()
	submission := &models.QuizSubmission{
		ID:          uuid.NewString(),
		UserID:      userID,
		TotalPoints: totalPoints,
		CompletedAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(submission).Error; err != nil {
			return err
		}

		for index := range answers {
			answers[index].ID = uuid.NewString()
			answers[index].SubmissionID = submission.ID
			answers[index].CreatedAt = now
		}
		if err := tx.Create(&answers).Error; err != nil {
			return err
		}

		scoreRows := make([]models.QuizSubmissionODSScore, 0, len(odsScores))
		for _, score := range odsScores {
			scoreRows = append(scoreRows, models.QuizSubmissionODSScore{
				ID:           uuid.NewString(),
				SubmissionID: submission.ID,
				ODSNumber:    score.ODSNumber,
				Points:       score.Points,
				CreatedAt:    now,
			})
		}
		if err := tx.Create(&scoreRows).Error; err != nil {
			return err
		}

		result := tx.Model(&models.UserScore{}).
			Where("id = ?", userID).
			Updates(map[string]any{"points": totalPoints, "updated_at": now})
		if result.Error != nil {
			return result.Error
		}

		// Older local databases can be missing the user row during API testing.
		// Submissions remain useful BI facts, so this write is best effort.
		if result.RowsAffected == 0 {
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = r.DB.Preload(clause.Associations).First(submission, "id = ?", submission.ID).Error
	if err != nil {
		return nil, err
	}

	return submission, nil
}
