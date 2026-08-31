package models

import "time"

type QuizQuestion struct {
	ID        string             `gorm:"primaryKey;type:uuid" json:"id"`
	Number    int                `gorm:"uniqueIndex;not null" json:"number"`
	Text      string             `gorm:"type:text;not null" json:"text"`
	ODSNumber int                `gorm:"not null;check:ods_number >= 1 AND ods_number <= 17" json:"ods_number"`
	Active    bool               `gorm:"not null;default:true" json:"active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Options   []QuizAnswerOption `gorm:"foreignKey:QuestionID" json:"options"`
}

func (QuizQuestion) TableName() string {
	return "quiz_questions"
}

type QuizAnswerOption struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	QuestionID  string    `gorm:"type:uuid;not null;uniqueIndex:idx_quiz_option_question_value" json:"question_id"`
	OptionValue int       `gorm:"not null;uniqueIndex:idx_quiz_option_question_value;check:option_value >= 1 AND option_value <= 5" json:"option_value"`
	Label       string    `gorm:"type:text;not null" json:"label"`
	Score       int       `gorm:"not null;default:0" json:"score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (QuizAnswerOption) TableName() string {
	return "quiz_answer_options"
}

type QuizSubmission struct {
	ID          string                   `gorm:"primaryKey;type:uuid" json:"id"`
	UserID      string                   `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalPoints int                      `gorm:"not null;default:0" json:"total_points"`
	CompletedAt time.Time                `gorm:"not null;index" json:"completed_at"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Answers     []QuizSubmissionAnswer   `gorm:"foreignKey:SubmissionID" json:"answers,omitempty"`
	ODSScores   []QuizSubmissionODSScore `gorm:"foreignKey:SubmissionID" json:"ods_scores,omitempty"`
}

func (QuizSubmission) TableName() string {
	return "quiz_submissions"
}

type QuizSubmissionAnswer struct {
	ID             string    `gorm:"primaryKey;type:uuid" json:"id"`
	SubmissionID   string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_quiz_submission_question" json:"submission_id"`
	QuestionID     string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_quiz_submission_question" json:"question_id"`
	AnswerOptionID string    `gorm:"type:uuid;not null;index" json:"answer_option_id"`
	QuestionNumber int       `gorm:"not null;index" json:"question_number"`
	ODSNumber      int       `gorm:"not null;index" json:"ods_number"`
	OptionValue    int       `gorm:"not null" json:"option_value"`
	Points         int       `gorm:"not null;default:0" json:"points"`
	CreatedAt      time.Time `json:"created_at"`
}

func (QuizSubmissionAnswer) TableName() string {
	return "quiz_submission_answers"
}

type QuizSubmissionODSScore struct {
	ID           string    `gorm:"primaryKey;type:uuid" json:"id"`
	SubmissionID string    `gorm:"type:uuid;not null;index;uniqueIndex:idx_quiz_submission_ods" json:"submission_id"`
	ODSNumber    int       `gorm:"not null;index;uniqueIndex:idx_quiz_submission_ods" json:"ods_number"`
	Points       int       `gorm:"not null;default:0" json:"points"`
	CreatedAt    time.Time `json:"created_at"`
}

func (QuizSubmissionODSScore) TableName() string {
	return "quiz_submission_ods_scores"
}

type UserScore struct {
	ID     string `gorm:"primaryKey;type:uuid"`
	Points int
}

func (UserScore) TableName() string {
	return "users"
}

type SubmitQuizRequest struct {
	Answers []SubmittedAnswer `json:"answers"`
}

type SubmittedAnswer struct {
	QuestionID     string `json:"question_id"`
	AnswerOptionID string `json:"answer_option_id"`
	QuestionNumber int    `json:"question_number"`
	OptionValue    int    `json:"option_value"`
}

type ODSScore struct {
	ODSNumber int `json:"ods_number"`
	Points    int `json:"points"`
}

type SubmitQuizResponse struct {
	SubmissionID string     `json:"submission_id"`
	UserID       string     `json:"user_id"`
	TotalPoints  int        `json:"total_points"`
	ODSScores    []ODSScore `json:"ods_scores"`
}
