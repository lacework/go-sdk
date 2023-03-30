use encoding_rs::*;

#[no_mangle]
pub extern fn run() -> () {

    unsafe { 
        hello();
        abort(0, 0, 0, 0);
    }

    let s = String::from("Hello from Rust!");
    log(&s);
    unsafe { lacework_api(0, 0); }
    () 
}

fn log(s: &str) {
    unsafe { http_request(0, 0); }
}

#[link(wasm_import_module = "")]
extern {
    fn hello();

    fn abort(a: u32, b: u32, c: u32, d: u32);

    #[link_name = "httpRequest"]
    fn http_request(a: u64, b: u64);

    // you can also explicitly specify the name to import, this imports `bar`
    // instead of `baz` from `the-wasm-import-module`.
    #[link_name = "laceworkAPI"]
    fn lacework_api(a: u64, b: u64);
}
