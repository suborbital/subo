use suborbital::runnable;

struct Example{}

impl runnable::Runnable for Example {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let in_string = String::from_utf8(input).unwrap();
    
        Some(String::from(format!("hello {}", in_string)).as_bytes().to_vec())
    }
}


// initialize the runner, do not edit below //
static EXAMPLE_RUNNER: &Example = &Example{};

#[no_mangle]
pub extern fn init() {
    runnable::set(EXAMPLE_RUNNER);
}
