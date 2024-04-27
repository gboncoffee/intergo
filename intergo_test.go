package intergo

import "testing"

func initTestingContext() InterContext {
	var ctx InterContext
	ctx.Init()
	ctx.AddLocale("pt_BR.UTF-8", map[string]string{"hello": "olá"})
	return ctx
}

func TestAddLocale(t *testing.T) {
	ctx := initTestingContext()
	if ctx.languages["pt"] == nil {
		t.Fatalf("Map entry for pt is nil")
	}
	if ctx.languages["pt"]["BR"] == nil {
		t.Fatalf("Map entry for pt_BR is nil")
	}
	if ctx.languages["pt"]["BR"]["hello"] != "olá" {
		t.Fatalf("Entry for hello is %v", ctx.languages["pt"]["BR"]["hello"])
	}
	ctx.AddLocale("eo_IN.UTF-8", map[string]string{"hello": "saluton"})
	if ctx.languages["eo"] == nil {
		t.Fatalf("Map entry for eo is nil")
	}
	if ctx.languages["eo"]["IN"] == nil {
		t.Fatalf("Map entry for pt_BR is nil")
	}
	if ctx.languages["eo"]["IN"]["hello"] != "saluton" {
		t.Fatalf("Entry for hello is %v", ctx.languages["eo"]["IN"]["hello"])
	}
}

func TestGetLocale(t *testing.T) {
	ctx := initTestingContext()
	txt, err := ctx.GetFromLocale("hello", "pt_BR.UTF-8")
	if err != nil {
		t.Fatalf("Got error on GetFromLocale hello pt_BR: %v", err)
	}
	if txt != "olá" {
		t.Fatalf("Text for hello pt_BR is wrong: %v", txt)
	}
	txt, err = ctx.GetFromLocale("hello", "en_US.UTF-8")
	if err != nil {
		t.Fatalf("Got error getting unset locale en_US: %v", err)
	}
	if txt != "hello" {
		t.Fatalf("Got wrong text getting unset locale en_US: %v", txt)
	}
	txt, err = ctx.GetFromLocale("hello", "pt_PT.UTF-8")
	if err != nil {
		t.Fatalf("Got error getting unset locale pt_PT: %v", err)
	}
	if txt != "olá" {
		t.Fatalf("Got wrong text getting unset locale pt_PT: %v", txt)
	}
}

func TestSetLocale(t *testing.T) {
	ctx := initTestingContext()
	ctx.SetPreferedLocale("pt_BR")
	txt := ctx.Get("hello")
	if txt != "olá" {
		t.Fatalf("Text for hello in prefered locale pt_BR is wrong: %v", txt)
	}
	ctx.SetPreferedLocale("pt_PT")
	txt = ctx.Get("hello")
	if txt != "olá" {
		t.Fatalf("Text for hello in prefered locale pt_PT is wrong: %v", txt)
	}
}
