package utils

import (
	"regexp"
	"strings"
)

//! CreateSlug génère un slug URL-friendly depuis un texte
func CreateSlug(text string) string {
	//! Convertir en minuscule
	slug := strings.ToLower(text)

	//! Remplacer espaces par tirets
	slug = strings.ReplaceAll(slug, " ", "-")

	//! Garder seulement lettres, chiffres, tirets
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	//! Supprimer tirets multiples
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	//! Trim tirets début/fin
	slug = strings.Trim(slug, "-")

	return slug
}

//! SplitFullName sépare un nom complet en prénom et nom
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
