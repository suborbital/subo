use crate::net::fetch;

#[no_mangle]
pub fn run(input: Vec<u8>) -> Option<Vec<u8>> {
    let url_string = String::from_utf8(input).unwrap();

    let data = fetch(url_string.as_str());

    Some(data)
}