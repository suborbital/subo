package template

import "fmt"

// ForLang returns a template for the lang provided
func ForLang(lang string) (string, string, error) {
	switch lang {
	case "rust":
		return "run.rs", rustTmpl(), nil
	case "swift":
		return "run.swift", swiftTmpl(), nil
	default:
		return "", "", fmt.Errorf("no template available for lang: %q", lang)
	}
}

func rustTmpl() string {
	tmpl := `
pub fn run(input: String) -> Option<String> {
	
	let out = String::from(format!("hello {}", input));
	
	return Some(out);
}`
	return tmpl
}

func swiftTmpl() string {
	tmpl := `
func run(input: String) -> String {

	return "hello " + input
}`

	return tmpl
}
