use suborbital::{req, runnable};

struct Helloworld{}

impl runnable::Runnable for Helloworld {
    fn run(&self, _: Vec<u8>) -> Option<Vec<u8>> {
        let body = req::body_raw();
        let body_string = util::to_string(body);
    
        Some(String::from(format!("hello {}", body_string)).as_bytes().to_vec())
    }
}


// initialize the runner, do not edit below //
static RUNNABLE: &Helloworld = &Helloworld{};

#[no_mangle]
pub extern fn init() {
    runnable::set(RUNNABLE);
}
