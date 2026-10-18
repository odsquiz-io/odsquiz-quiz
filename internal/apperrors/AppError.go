package apperrors

import "errors"

type Kind string
type Code string

const (
	KindBadRequest   Kind = "bad_request"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindInternal     Kind = "internal"
)

const (
	CodeInvalidRequest     Code = "invalid_request"
	CodeUnauthorized       Code = "unauthorized"
	CodeQuestionNotFound   Code = "question_not_found"
	CodeAnswerNotFound     Code = "answer_not_found"
	CodeIncompleteQuiz     Code = "incomplete_quiz"
	CodeDuplicateQuestion  Code = "duplicate_question"
	CodeUserIDUnavailable  Code = "user_id_unavailable"
	CodeInvalidCredentials Code = "invalid_credentials"
	CodeAccountSuspended   Code = "account_suspended"
	CodeInternal           Code = "internal_error"
)

type AppError struct {
	Kind Kind
	Code Code
	Err  error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return string(e.Code)
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func BadRequest(code Code, err error) error {
	return &AppError{Kind: KindBadRequest, Code: code, Err: err}
}

func Unauthorized(code Code, err error) error {
	return &AppError{Kind: KindUnauthorized, Code: code, Err: err}
}

func Forbidden(code Code, err error) error {
	return &AppError{Kind: KindForbidden, Code: code, Err: err}
}

func NotFound(code Code, err error) error {
	return &AppError{Kind: KindNotFound, Code: code, Err: err}
}

func Internal(err error) error {
	return &AppError{Kind: KindInternal, Code: CodeInternal, Err: err}
}

func From(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
