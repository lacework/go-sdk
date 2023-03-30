use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern {
    pub fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet(name: &str) {
    alert(&format!("Hello, {}!", name));
}

#[wasm_bindgen]
pub fn abort(a: u32, b: u32, c: u32, d: u32) {}

#[wasm_bindgen]
pub fn http_request(a: u64, b: u64) {}

#[wasm_bindgen]
pub fn lacework_api(a: u64, b: u64) {}
