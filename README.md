# InterGo

Small library for simple internationalization in Go.

## Usage

```go
var ctx InterContext
ctx.Init()
ctx.AddLocale("pt_BR", map[string]string{"hello": "olá"})
```

Optionally set the prefered locale to properly use `ctx.Get()`:
```go
ctx.SetPreferedLocale("pt_BR")
```

### Get localized strings:

This returns "olá", as we have a "pt_BR" locale set.
```go
txt, err := ctx.GetFromLocale("hello", "pt_BR")
```

This returns "olá", as we haven't set any Portuguese Portuguese locale, so it
falls back to other locales in the same language:
```go
txt, err = ctx.GetFromLocale("hello", "pt_PT")
```

This returns "hello", as we haven't set any English language locale, so it'll
just return the string we have passed.
```go
txt, err = ctx.GetFromLocale("hello", "en_US")
```

### Prefered locale.

It's possible to set a prefered locale. This way, we simply use `ctx.Get()`
to retrive strings instead of passing the locale every time:

```go
err := ctx.SetPreferedLocale(locale)
if err != nil {
    return fmt.Errorf("error parsing locale string: %v", locale)
}
txt := ctx.Get("hello")
```

Note how `ctx.Get()` does not need to return any error as it does not parses
a locale string.

It's also possible to automatically set the prefered locale from the
environment variables `LC_ALL` and `LANG`:
```go
err := ctx.AutoSetPreferedLocale()
if err != nil {
    return fmt.Errorf("error parsing environment variables: %v", err)
}
```
