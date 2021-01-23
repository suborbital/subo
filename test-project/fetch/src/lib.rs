use suborbital::runnable;
use suborbital::http;

struct Fetch{}

impl runnable::Runnable for Fetch {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let in_string = String::from_utf8(input).unwrap();
    
        let result = http::get(in_string.as_str());
        let result_string = String::from_utf8(result).unwrap();
    
        Some(result_string.as_bytes().to_vec())
    }
}


// initialize the runner, do not edit below //
static FETCH_RUNNER: &Fetch = &Fetch{};

#[no_mangle]
pub extern fn init() {
    runnable::set(FETCH_RUNNER);
}
