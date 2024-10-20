package i18n

import "embed"

//go:embed en
//go:embed pt

// LocaleFS is the filesystem containing the locale files.
var LocaleFS embed.FS

// SupportedLocales is the list of supported locales.
var SupportedLocales = map[string]bool{
	"en": true,
	"pt": true,
}
