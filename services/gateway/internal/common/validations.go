package common

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ---------- Variables ----------
var (
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{1,50}$`)
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ---------- FUNCTIONS ----------

func ValidatePositiveInt(id int64) (int64, error) {
	if id <= 0 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil
}

func ValidatePositiveIntString(id string) (int64, error) {
	if id == "" {
		return 0, errors.New("id is required")
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return ValidatePositiveInt(idInt)
}

func ValidateIDList(idList []int64) error {
	for _, id := range idList {
		if id <= 0 {
			return errors.New("id must be a positive integer")
		}
	}
	return nil
}

func ValidateIDListFromString(idList string) ([]int64, error) {
	idStrs := strings.Split(idList, ",")
	var ids []int64
	for _, s := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := ValidateIDList(ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func IsValidPassword(password string) bool {
	// Mínimo 8 caracteres
	if len(password) < 8 {
		return false
	}
	// Al menos una mayúscula
	hasUppercase, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUppercase {
		return false
	}
	// Al menos un caracter especial
	hasSpecial, _ := regexp.MatchString(`[^a-zA-Z0-9]`, password)
	if !hasSpecial {
		return false
	}
	return true
}
