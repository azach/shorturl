package storage

import (
	"github.com/sirupsen/logrus"
	"time"
)

type HitRange int

const (
	AllTime HitRange = iota
	Daily
	Weekly
	Minute
)

type Storage interface {
	Get(key string) (value string, exists bool)
	Set(key string, value string) error
	Hit(key string, viewedAt time.Time)
	GetHits(key string, asOf time.Time, hitRange HitRange) (int64, error)
}

func toPrecision(hitRange HitRange) int64 {
	switch hitRange {
	case AllTime:
		return 0
	case Daily:
		return 86400
	case Weekly:
		return 604800
	case Minute:
		return 60
	default:
		logrus.Errorf("Unknown range defined %v", hitRange)
		return 0
	}
}
