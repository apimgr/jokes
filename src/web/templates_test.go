package web

import "testing"

func TestInitTemplates(t *testing.T) {
	Templates = nil

	if err := InitTemplates(); err != nil {
		t.Fatalf("InitTemplates() error = %v", err)
	}
	if Templates == nil {
		t.Fatal("InitTemplates() left Templates nil")
	}
	if tpl := Templates.Lookup("content"); tpl == nil {
		t.Error("Templates missing expected \"content\" template from index.html")
	}
	if tpl := Templates.Lookup("header"); tpl == nil {
		t.Error("Templates missing expected \"header\" template")
	}
}
