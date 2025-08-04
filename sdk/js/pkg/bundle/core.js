import * as wasm from "./core_bg.wasm";
export * from "./core_bg.js";
import { __wbg_set_wasm } from "./core_bg.js";
__wbg_set_wasm(wasm);
wasm.__wbindgen_start();
