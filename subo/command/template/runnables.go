package template

import "fmt"

// ForLang returns a template for the lang provided
func ForLang(lang string) (string, string, error) {
	switch lang {
	case "rust":
		return "run.rs", rustTmpl(), nil
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
