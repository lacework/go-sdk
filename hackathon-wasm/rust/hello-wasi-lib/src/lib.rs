#[no_mangle]
pub extern "C" fn run(a: u32, b: u32, c: u32, d: u32) -> (u32) {
    println!("Hello, world!");
    (0)
}

#[link(wasm_import_module = "the-wasm-import-module")]
extern "C" {
    // imports the name `foo` from `the-wasm-import-module`
    fn hello();

    fn abort(a: u32, b: u32, c: u32, d: u32);

    fn httpRequest(a: u64, b: u64);

    // you can also explicitly specify the name to import, this imports `bar`
    // instead of `baz` from `the-wasm-import-module`.
    #[link_name = "laceworkAPI"]
    fn lacework_api(a: u64, b: u64);
}
