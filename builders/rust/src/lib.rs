/**
 * 
 * This file represents the Rust "API" for Hive WASM runnables. The functions defined herein are used to exchange data
 * between the host (Hive, written in Go) and the Runnable (a WASM module, in this case written in Rust). The Runnable's 
 * public facing `run` function does not need to concern itself with this API, and simply needs to have the this signature:
 * 
 * #[no_mangle]
 * pub fn run(input: Vec<u8>) -> Option<Vec<u8>>
 * 
 */

use std::mem;
use std::slice;

mod run;

extern {
    fn return_result(result_pointer: *const u8, result_size: i32, ident: i32);
    fn fetch(url_pointer: *const u8, url_size: i32, dest_pointer: *const u8, dest_max_size: i32, ident: i32) -> i32;
    fn print(pointer: *const u8, result_size: i32, ident: i32);
}

// keep track of the current identifier to use for calls across FFI
static mut CURRENT_IDENT: i32 = 0;

#[no_mangle]
pub extern fn allocate(size: i32) -> *const u8 {
    let mut buffer = Vec::with_capacity(size as usize);
    let buffer_slice = buffer.as_mut_slice();
    let pointer = buffer_slice.as_mut_ptr();
    mem::forget(buffer_slice);

    pointer as *const u8
}

#[no_mangle]
pub extern fn deallocate(pointer: *const u8, size: i32) {
    unsafe {
        let _ = slice::from_raw_parts(pointer, size as usize);
    }
}

#[no_mangle]
pub extern fn run_e(pointer: *const u8, size: i32, ident: i32) {
    unsafe { CURRENT_IDENT = ident };

    // rebuild the memory into something usable
    let in_slice: &[u8] = unsafe { 
        slice::from_raw_parts(pointer, size as usize) 
    };

    let in_bytes = Vec::from(in_slice);

    // call the runnable and check its result
    let result: Vec<u8> = match run::run(in_bytes) {
        Some(val) => val,
        None => Vec::from("run returned no data"), 
    };

    let result_slice = result.as_slice();
    let result_size = result_slice.len();


    // call back to hive to return the result
    unsafe { 
        return_result(result_slice.as_ptr() as *const u8, result_size as i32, ident); 
    }
}

pub mod net {
    pub fn fetch(url: &str) -> Vec<u8> {
        let mut dest_pointer: *const u8;
        let mut dest_size: i32;
        let mut capacity: i32 = 256000;

        // make the request, and if the response size is greater than that of capacity, double the capacity and try again
        loop {
            let cap = &mut capacity;

            let mut dest_bytes = Vec::with_capacity(*cap as usize);
            let dest_slice = dest_bytes.as_mut_slice();
            dest_pointer = dest_slice.as_mut_ptr() as *const u8;
    
            // do the request over FFI
            dest_size = unsafe { super::fetch(url.as_ptr(), url.len() as i32, dest_pointer, *cap, super::CURRENT_IDENT) };

            if dest_size < 0 {
                return Vec::from(format!("request_failed:{}", dest_size))
            } else if dest_size > *cap {
                super::fmt::print(format!("doubling capacity, need {}", dest_size).as_str());
                *cap *= 2;
            } else {
                break;
            }
        }

        let result: &[u8] = unsafe {
            super::slice::from_raw_parts(dest_pointer, dest_size as usize)
        };

        return Vec::from(result)
    }
}

pub mod fmt {
    pub fn print(msg: &str) {
        let msg_vec = Vec::from(msg);
        let msg_slice = msg_vec.as_slice();
        let pointer = msg_slice.as_ptr();

        unsafe {super::print(pointer, msg_slice.len() as i32, super::CURRENT_IDENT) };
    }
}