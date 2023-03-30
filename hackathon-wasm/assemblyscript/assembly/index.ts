// The entry file of your WebAssembly module.
import {logging, httpRequest, writeFile} from './env'

function log(msg: string): void {
  let buf = String.UTF8.encode(msg);
  let ptr = changetype<usize>(buf);
  let len = buf.byteLength;

  logging(ptr, len);
}

class Raw {
  constructor(ptr: u64, len: u64) {
    this.ptr = ptr;
    this.len = len;
  }
  ptr: u64;
  len: u64;
}

function convert(buf: ArrayBuffer): Raw {
  const raw = new Raw(changetype<usize>(buf), buf.byteLength);
  return raw;
}

function HTTP(verb: string, url: string, headers: string, body: string): ArrayBuffer {
  let response = new ArrayBuffer(1000000);

  let payload = String.UTF8.encode(`{
    "verb": "${verb}",
    "url":  "${url}",
    "headers": ${headers},
    "body": ${body},
    "response": ${changetype<usize>(response)}
  }`);

  const raw = convert(payload)

  httpRequest(raw.ptr, raw.len);

  return response;
}

function file(path: string, contentPtr: u64, contentLen: u64): void {
  const buf = String.UTF8.encode(path)
  const raw = convert(buf)

  // TODO: hack to ensure that functions are properly matched
  writeFile(raw.ptr, raw.len, contentPtr, contentLen);
}

const token = ""

export function chat(ptr: u64, length: i32): void {
  log("chat");

  let r = new Uint8Array(length);

  for (let i = 0; i < length; i++) {
    r[i] = load<u8>(usize(ptr + i));
  }

  const content = String.UTF8.decode(r.buffer);

  const resp = HTTP(
     "POST",
     "https://api.openai.com/v1/chat/completions",
     `{"Authorization": "Bearer ${token}", "Content-Type": "application/json"}`,
     `{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "${content}"}]}`)


  const cliMsg = String.UTF8.decode(resp).split(',').filter(x => x.includes("content"))[0].replace('"content":"', "").replace('"}', '');

  log(cliMsg);

  return;
}

export function present(): void {
  log("present");

  const resp = HTTP("GET", "https://i.imgflip.com/7ggtv4.jpg", "{}", "{}");

  file("/Users/j0n/Desktop/wasm.png", changetype<usize>(resp), resp.byteLength);
}
