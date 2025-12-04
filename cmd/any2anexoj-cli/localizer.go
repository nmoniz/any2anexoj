package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed translations/*.json
var translationsFS embed.FS

type Localizer struct {
	*i18n.Localizer
}

func NewLocalizer(lang string) (*Localizer, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFileFS(translationsFS, "translations/en.json")
	if err != nil {
		return nil, fmt.Errorf("loading english messages: %w", err)
	}

	_, err = bundle.LoadMessageFileFS(translationsFS, "translations/pt.json")
	if err != nil {
		return nil, fmt.Errorf("loading portuguese messages: %w", err)
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	return &Localizer{
		Localizer: localizer,
	}, nil
}

func (t Localizer) Translate(key string, count int, values map[string]any) string {
	txt, err := t.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: values,
		PluralCount:  count,
	})
	if err != nil {
		slog.Error("failed to translate message", slog.Any("err", err))
		return "<ERROR>"
	}
	return txt
}
