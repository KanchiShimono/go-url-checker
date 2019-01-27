package repository

import (
	"net/url"
	"time"
)

// FileRepositoryReader read test condition of http health checke.
type FileRepositoryReader interface {
	ReadAll() ([]ConditionHTTPCheck, error)
}

// FileRepositoryWriter wirte result of http health checke.
type FileRepositoryWriter interface {
	WriteAll([]ResultHTTPCheck) error
}

// ConditionHTTPCheck is input data for http check.
type ConditionHTTPCheck struct {
	URL         *url.URL
	StatusCode  int
	Timeout     time.Duration
	Description string
}

// ResultHTTPCheck is result data for http check.
type ResultHTTPCheck struct {
	TimeStamp   string
	URL         *url.URL
	Result      error
	Description string
}
