package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CreateSlug génère un slug URL-friendly depuis un texte
func CreateSlug(text string) string {
	// 1. Convertir en minuscule
	slug := strings.ToLower(text)

	// 2. Supprimer les accents
	slug = removeAccents(slug)

	// 3. Remplacer espaces et underscores par tirets
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// 4. Garder seulement lettres, chiffres, tirets
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// 5. Supprimer tirets multiples
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// 6. Trim tirets début/fin
	slug = strings.Trim(slug, "-")

	return slug
}

// removeAccents supprime les accents d'un texte
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// SplitFullName sépare un nom complet en prénom et nom
func SplitFullName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}

	firstName = parts[0]
	if len(parts) > 1 {
		lastName = strings.Join(parts[1:], " ")
	}

	return firstName, lastName
}
