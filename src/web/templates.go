package web

import (
	"embed"
	"html/template"
)

//go:embed templates/* static/*
var EmbeddedFiles embed.FS

var Templates *template.Template

func InitTemplates() error {
	var err error
	Templates, err = template.ParseFS(EmbeddedFiles, "templates/*.html")
	return err
}
