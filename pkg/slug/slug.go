package slug

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	slugUnsafeChars = regexp.MustCompile(`[^a-z0-9 ]+`)
	slugSpaces      = regexp.MustCompile(` +`)
)

func ToSlug(s string) string {
	identifier := uuid.New().String()[:6]

	var title string

	if len(s) > 18 {
		title = s[:18]
	} else {
		title = s
	}

	baseSlug := strings.ToLower(title)
	baseSlug = slugUnsafeChars.ReplaceAllString(baseSlug, "")
	baseSlug = slugSpaces.ReplaceAllString(baseSlug, "-")
	baseSlug = strings.Trim(baseSlug, "-")

	if baseSlug == "" {
		return identifier
	}

	return baseSlug + "-" + identifier
}
