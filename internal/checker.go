package internal

import (
	"errors"
	"forum/pkg/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func CheckPost(p *models.Post) int {

	if len(RemoveSpace(p.Title)) == 0 || len(p.Content) == 0 {
		return 1
	}

	if len(p.Content) > 250 {
		return 2
	}

	if len(p.Title) > 60 {
		return 3
	}

	return 0
}

func IsValidCtgry(r *http.Request) ([]int, error) {
	categories := r.Form["category"]

	var categoryInt []int
	for _, el := range categories {
		categoryID, err := strconv.Atoi(string(el))
		if err != nil || categoryID > 5 {
			return []int{}, errors.New("Invalid category ID")
		}

		categoryInt = append(categoryInt, categoryID)

	}
	return categoryInt, nil
}

func CheckName(s string) bool {
	if len(RemoveSpace(s)) == 0 || len(s) > 50 {
		return false
	}
	loginConvention := "^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$"
	if re, _ := regexp.Compile(loginConvention); !re.MatchString(s) {
		return false
	}
	return true
}
func CheckEmail(email string) bool {
	if len(RemoveSpace(email)) == 0 || len(email) > 320 {
		return false
	}

	// Valid characters for the email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return false
	}
	return true
}

func RemoveSpace(s string) string {
	return strings.Trim(s, " ")
}
