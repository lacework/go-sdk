// The entry file of your WebAssembly module.
import {hello, httpRequest, laceworkAPI } from './env'


function log(msg: string): void {
  let buf = String.UTF8.encode(msg);
  let ptr = changetype<usize>(buf);
  let len = buf.byteLength;

  httpRequest(ptr, len);
}

export function run(): void {

  hello();

  log("https://google.com")

  let result_buf = new ArrayBuffer(2048);
  const data = `{"function": "my-lacework-api-function", "arg1": "my arguments", "result_ptr": ${changetype<usize>(result_buf)}, "result_len": ${result_buf.byteLength}}`;

  let arg2 = String.UTF8.encode(data);

  laceworkAPI(changetype<usize>(arg2), arg2.byteLength);

  const response = String.UTF8.decode(result_buf);

  log(response);

  return;
}
