package utils

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func AtoiSafe(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
