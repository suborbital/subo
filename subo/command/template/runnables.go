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
#[no_mangle]
pub fn run(input: Vec<u8>) -> Option<Vec<u8>> {
	let in_string = String::from_utf8(input).unwrap();

	Some(String::from(format!("hello {}", in_string)).as_bytes().to_vec())
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
