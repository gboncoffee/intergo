// Package intergo implements a simple library for internationalized strings.
// The library manages a hash map in the form `map["language"]["locale"]`.
// The supported format for locale strings is `language_locale.encoding`. The
// encoding part is actually ignored, and the form `language_locale` also works.
// Of course, it's case-sensitive and the recommended form is `language_LOCALE`,
// e.g., `en_US` is a locale for American English, and `pt_BR` is for Brazilian
// Portuguese.
//
// ## Example usage:
//
// ```go
// var ctx InterContext
// ctx.Init()
// ctx.AddLocale("pt_BR", map[string]string{"hello": "olá"})
// ```
//
// Optionally set the prefered locale to properly use `ctx.Get()`:
// ```go
// ctx.SetPreferedLocale("pt_BR")
// ```
//
// ### Get localized strings:
//
// This returns "olá", as we have a "pt_BR" locale set.
// ```go
// txt, err := ctx.GetFromLocale("hello", "pt_BR")
// ```
//
// This returns "olá", as we haven't set any Portuguese Portuguese locale, so it
// falls back to other locales in the same language:
// ```go
// txt, err = ctx.GetFromLocale("hello", "pt_PT")
// ```
//
// This returns "hello", as we haven't set any English language locale, so it'll
// just return the string we have passed.
// ```go
// txt, err = ctx.GetFromLocale("hello", "en_US")
// ```
//
// ### Prefered locale.
//
// It's possible to set a prefered locale. This way, we simply use `ctx.Get()`
// to retrive strings instead of passing the locale every time:
//
// ```go
// err := ctx.SetPreferedLocale(locale)
// if err != nil {
//     return fmt.Errorf("error parsing locale string: %v", locale)
// }
// txt := ctx.Get("hello")
// ```
//
// Note how `ctx.Get()` does not need to return any error as it does not parses
// a locale string.
//
// It's also possible to automatically set the prefered locale from the
// environment variables `LC_ALL` and `LANG`:
// ```go
// err := ctx.AutoSetPreferedLocale()
// if err != nil {
//     return fmt.Errorf("error parsing environment variables: %v", err)
// }
// ```
package intergo

import (
	"fmt"
	"os"
)

// Returns the lang them the locale, just as the order in the string. E.g.,
// parsing "pt_BR.UTF-8" will return ("pt", "BR", nil). Works also without the
// encoding specification. We don't "support" nothing besides UTF-8 anyways.
func parseLocaleString(locale string) (string, string, error) {
	var lang string
	var local string

	n, err := fmt.Sscanf(locale, "%2s_%2s", &lang, &local)
	if err != nil {
		return "", "", err
	}
	if n != 2 {
		return "", "", fmt.Errorf("unparsable locale string %v", locale)
	}

	return lang, local, nil
}

// The type for a specific locale, i.e., the map with internationalized entries.
// E.g., the map `br` may have entries `br["hello"] == "olá"`.
type Locale map[string]string
// A collection of locales with the same language. E.g., `en_US` and `en_GB` are
// in the same Language map `en`.
type Language map[string]Locale
// The library context itself.
type InterContext struct {
	languages    map[string]Language
	prefered     Locale
	preferedLang Language
}

// Initializes the language map, should be called after instanciating an
// InterContext. Usually called right after the application startup.
func (ctx *InterContext) Init() {
	ctx.languages = make(map[string]Language)
}

// Automatically sets the prefered locale with the variables `LC_ALL` and
// `LANG`. Basically tries `LC_ALL`, and if it cannot parse a locale from it,
// tries `LC_LANG`.
func (ctx *InterContext) AutoSetPreferedLocale() error {
	lcvar := os.Getenv("LC_ALL")
	err := ctx.SetPreferedLocale(lcvar)
	if err == nil {
		return nil
	}
	lcvar = os.Getenv("LANG")
	return ctx.SetPreferedLocale(lcvar)
}

// Adds a new mapping of strings, i.e., a new locale, to the context. Usually
// called for all the supported locales right after the context initialization.
func (ctx *InterContext) AddLocale(locale string, entries map[string]string) error {
	lang, local, err := parseLocaleString(locale)
	if err != nil {
		return err
	}

	if ctx.languages[lang] == nil {
		ctx.languages[lang] = Language{local: entries}
	} else {
		ctx.languages[lang][local] = entries
	}

	return nil
}

// Sets the prefered locale.
func (ctx *InterContext) SetPreferedLocale(locale string) error {
	lang, local, err := parseLocaleString(locale)
	if err != nil {
		return err
	}

	ctx.preferedLang = ctx.languages[lang]
	if ctx.preferedLang != nil {
		ctx.prefered = ctx.preferedLang[local]
	}

	return nil
}

// Gets an internationalized string with the prefered locale.
func (ctx *InterContext) Get(text string) string {
	if ctx.preferedLang == nil {
		return text
	}
	var localTxt string
	if ctx.prefered != nil {
		localTxt = ctx.prefered[text]
		if localTxt != "" {
			return localTxt
		}
	}
	for _, l := range ctx.preferedLang {
		localTxt = l[text]
		if localTxt != "" {
			return localTxt
		}
	}

	return text
}

// Gets an internationalized string from a specific locale.
func (ctx *InterContext) GetFromLocale(text string, locale string) (string, error) {
	lang, local, err := parseLocaleString(locale)
	if err != nil {
		return text, err
	}

	langMap := ctx.languages[lang]
	if langMap == nil {
		return text, nil
	}

	localMap := langMap[local]
	if localMap == nil {
		for _, l := range langMap {
			localTxt := l[text]
			if localTxt != "" {
				return localTxt, nil
			}
		}
		return text, nil
	}

	localTxt := localMap[text]
	if localTxt == "" {
		return text, nil
	}

	return localTxt, nil
}
