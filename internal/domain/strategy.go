package domain

import (
	"time"

	"github.com/google/uuid"
)

type Strategy struct {
	ID          uuid.UUID
	Name        string
	Description string
	Language    StrategyLanguage
	Code        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

type StrategyLanguage int

const (
	StrategyLanguageGo StrategyLanguage = iota
	StrategyLanguagePython
)

func (sl StrategyLanguage) String() string {
	switch sl {
	case StrategyLanguageGo:
		return "GO"
	case StrategyLanguagePython:
		return "PYTHON"
	default:
		return "UNKNOWN"
	}
}

func (s Strategy) Validate() error {
	if s.Name == "" {
		return ErrInvalidStrategy{Reason: "name cannot be empty"}
	}
	if s.Code == "" {
		return ErrInvalidStrategy{Reason: "code cannot be empty"}
	}
	if s.Language != StrategyLanguagePython && s.Language != StrategyLanguageGo {
		return ErrInvalidStrategy{Reason: "unsupported language"}
	}

	return nil
}

type ErrInvalidStrategy struct {
	Reason string
}

func (e ErrInvalidStrategy) Error() string {
	return "invalid strategy: " + e.Reason
}
