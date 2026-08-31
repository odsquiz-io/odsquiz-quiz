package migrations

import (
	"context"
	"fmt"

	"github.com/kauanpecanha/odsquiz-quiz/internal/models"
	"github.com/kauanpecanha/odsquiz-quiz/internal/seed"
	"github.com/kauanpecanha/odsquiz-quiz/pkg/config"
	"github.com/kauanpecanha/odsquiz-quiz/pkg/database"
)

func RunMigrations(ctx context.Context, cfg *config.Config) error {
	db, err := database.NewPostgresConnection(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&models.QuizQuestion{},
		&models.QuizAnswerOption{},
		&models.QuizSubmission{},
		&models.QuizSubmissionAnswer{},
		&models.QuizSubmissionODSScore{},
	)
	if err != nil {
		return fmt.Errorf("migrate quiz tables: %w", err)
	}

	if err := seed.Questions(db); err != nil {
		return fmt.Errorf("seed quiz questions: %w", err)
	}

	return nil
}
