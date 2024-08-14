package models

import (
    "gorm.io/gorm"
    "github.com/go-playground/validator/v10"
)

type Strategy struct {
    gorm.Model
    UserID      uint   `json:"user_id" gorm:"index" validate:"required"`
    Name        string `json:"name" gorm:"not null" validate:"required,min=3,max=50"`
    Description string `json:"description" validate:"max=500"`
    Parameters  string `json:"parameters" validate:"required,json"`
}

func (s *Strategy) Validate() error {
    validate := validator.New()
    return validate.Struct(s)
}