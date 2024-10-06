package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	MinimumRequirementType     = "minimum"
	RecommendedRequirementType = "recommended"
	MaximumRequirementType     = "maximum"
	WindowsOSRequirement       = "windows"
	MacOSRequirement           = "mac"
	LinuxOSRequirement         = "linux"
)

type RequirementType struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Potential string `gorm:"not null;type:enum('minimum','recommended','maximum')" validate:"required,enum_potential"`
	OS        string `gorm:"not null;type:enum('windows','mac','linux')" validate:"required,enum_os"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (rt *RequirementType) ValidateRequirementType() error {
	Init()

	err := validate.Struct(rt)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
