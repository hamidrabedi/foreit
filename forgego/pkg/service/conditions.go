package service

import (
	"github.com/forgego/forge/pkg/models"
)

func newStringCondition(field, op string, value interface{}) models.Condition {
	return models.NewStringCondition(field, op, value)
}

func newIntCondition(field, op string, value interface{}) models.Condition {
	return models.NewIntCondition(field, op, value)
}

func newInCondition(field string, values interface{}) models.Condition {
	return models.NewInCondition(field, values)
}

func newIsNullCondition(field string, null bool) models.Condition {
	return models.NewIsNullCondition(field, null)
}

