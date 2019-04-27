package gpt

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ArrivalTime struct {
	time.Time
}

func (t *ArrivalTime) UnmarshalJSON(b []byte) error {
	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}
	layout := "2006-01-02T15:04:05"
	parsed, err := time.Parse(layout, dateString)
	t.Time = parsed
	return err
}
func (p ArrivalTime) Value() (driver.Value, error) {
	return driver.Value(p.Time), nil
}

type Date struct {
	time.Time
}

func (t *Date) UnmarshalJSON(b []byte) error {
	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}
	layout := "2006-01-02"
	parsed, err := time.Parse(layout, dateString)
	t.Time = parsed
	return err
}
func (p Date) Value() (driver.Value, error) {
	return driver.Value(p.Time), nil
}

type UpdateTime struct {
	time.Time
}

func (t *UpdateTime) UnmarshalJSON(b []byte) error {
	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	parsed, err := time.Parse(layout, dateString)
	t.Time = parsed
	return err
}

func (p UpdateTime) Value() (driver.Value, error) {
	return driver.Value(p.Time), nil
}

type SimpleTime struct {
	time.Time
}

func (t *SimpleTime) UnmarshalJSON(b []byte) error {
	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}
	layout := "15:04"
	parsed, err := time.Parse(layout, dateString)
	t.Time = parsed
	return err
}

func (p SimpleTime) Value() (driver.Value, error) {
	return driver.Value(p.Time), nil
}

type SimpleTimeWithSeconds struct {
	time.Time
}

func (t *SimpleTimeWithSeconds) UnmarshalJSON(b []byte) error {
	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}
	layout := "15:04:05"
	parsed, err := time.Parse(layout, dateString)
	t.Time = parsed
	return err
}

func (p SimpleTimeWithSeconds) Value() (driver.Value, error) {
	return driver.Value(p.Time), nil
}
