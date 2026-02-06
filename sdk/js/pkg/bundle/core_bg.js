let wasm;
export function __wbg_set_wasm(val) {
    wasm = val;
}

function addToExternrefTable0(obj) {
    const idx = wasm.__externref_table_alloc();
    wasm.__wbindgen_externrefs.set(idx, obj);
    return idx;
}

function _assertClass(instance, klass) {
    if (!(instance instanceof klass)) {
        throw new Error(`expected instance of ${klass.name}`);
    }
}

function debugString(val) {
    // primitive types
    const type = typeof val;
    if (type == 'number' || type == 'boolean' || val == null) {
        return  `${val}`;
    }
    if (type == 'string') {
        return `"${val}"`;
    }
    if (type == 'symbol') {
        const description = val.description;
        if (description == null) {
            return 'Symbol';
        } else {
            return `Symbol(${description})`;
        }
    }
    if (type == 'function') {
        const name = val.name;
        if (typeof name == 'string' && name.length > 0) {
            return `Function(${name})`;
        } else {
            return 'Function';
        }
    }
    // objects
    if (Array.isArray(val)) {
        const length = val.length;
        let debug = '[';
        if (length > 0) {
            debug += debugString(val[0]);
        }
        for(let i = 1; i < length; i++) {
            debug += ', ' + debugString(val[i]);
        }
        debug += ']';
        return debug;
    }
    // Test for built-in
    const builtInMatches = /\[object ([^\]]+)\]/.exec(toString.call(val));
    let className;
    if (builtInMatches && builtInMatches.length > 1) {
        className = builtInMatches[1];
    } else {
        // Failed to match the standard '[object ClassName]'
        return toString.call(val);
    }
    if (className == 'Object') {
        // we're a user defined class or Object
        // JSON.stringify avoids problems with cycles, and is generally much
        // easier than looping through ownProperties of `val`.
        try {
            return 'Object(' + JSON.stringify(val) + ')';
        } catch (_) {
            return 'Object';
        }
    }
    // errors
    if (val instanceof Error) {
        return `${val.name}: ${val.message}\n${val.stack}`;
    }
    // TODO we could test for more things here, like `Set`s and `Map`s.
    return className;
}

function getArrayU8FromWasm0(ptr, len) {
    ptr = ptr >>> 0;
    return getUint8ArrayMemory0().subarray(ptr / 1, ptr / 1 + len);
}

let cachedDataViewMemory0 = null;
function getDataViewMemory0() {
    if (cachedDataViewMemory0 === null || cachedDataViewMemory0.buffer.detached === true || (cachedDataViewMemory0.buffer.detached === undefined && cachedDataViewMemory0.buffer !== wasm.memory.buffer)) {
        cachedDataViewMemory0 = new DataView(wasm.memory.buffer);
    }
    return cachedDataViewMemory0;
}

function getStringFromWasm0(ptr, len) {
    ptr = ptr >>> 0;
    return decodeText(ptr, len);
}

let cachedUint8ArrayMemory0 = null;
function getUint8ArrayMemory0() {
    if (cachedUint8ArrayMemory0 === null || cachedUint8ArrayMemory0.byteLength === 0) {
        cachedUint8ArrayMemory0 = new Uint8Array(wasm.memory.buffer);
    }
    return cachedUint8ArrayMemory0;
}

function handleError(f, args) {
    try {
        return f.apply(this, args);
    } catch (e) {
        const idx = addToExternrefTable0(e);
        wasm.__wbindgen_exn_store(idx);
    }
}

function isLikeNone(x) {
    return x === undefined || x === null;
}

function passArray8ToWasm0(arg, malloc) {
    const ptr = malloc(arg.length * 1, 1) >>> 0;
    getUint8ArrayMemory0().set(arg, ptr / 1);
    WASM_VECTOR_LEN = arg.length;
    return ptr;
}

function passStringToWasm0(arg, malloc, realloc) {
    if (realloc === undefined) {
        const buf = cachedTextEncoder.encode(arg);
        const ptr = malloc(buf.length, 1) >>> 0;
        getUint8ArrayMemory0().subarray(ptr, ptr + buf.length).set(buf);
        WASM_VECTOR_LEN = buf.length;
        return ptr;
    }

    let len = arg.length;
    let ptr = malloc(len, 1) >>> 0;

    const mem = getUint8ArrayMemory0();

    let offset = 0;

    for (; offset < len; offset++) {
        const code = arg.charCodeAt(offset);
        if (code > 0x7F) break;
        mem[ptr + offset] = code;
    }
    if (offset !== len) {
        if (offset !== 0) {
            arg = arg.slice(offset);
        }
        ptr = realloc(ptr, len, len = offset + arg.length * 3, 1) >>> 0;
        const view = getUint8ArrayMemory0().subarray(ptr + offset, ptr + len);
        const ret = cachedTextEncoder.encodeInto(arg, view);

        offset += ret.written;
        ptr = realloc(ptr, len, offset, 1) >>> 0;
    }

    WASM_VECTOR_LEN = offset;
    return ptr;
}

function takeFromExternrefTable0(idx) {
    const value = wasm.__wbindgen_externrefs.get(idx);
    wasm.__externref_table_dealloc(idx);
    return value;
}

let cachedTextDecoder = new TextDecoder('utf-8', { ignoreBOM: true, fatal: true });
cachedTextDecoder.decode();
const MAX_SAFARI_DECODE_BYTES = 2146435072;
let numBytesDecoded = 0;
function decodeText(ptr, len) {
    numBytesDecoded += len;
    if (numBytesDecoded >= MAX_SAFARI_DECODE_BYTES) {
        cachedTextDecoder = new TextDecoder('utf-8', { ignoreBOM: true, fatal: true });
        cachedTextDecoder.decode();
        numBytesDecoded = len;
    }
    return cachedTextDecoder.decode(getUint8ArrayMemory0().subarray(ptr, ptr + len));
}

const cachedTextEncoder = new TextEncoder();

if (!('encodeInto' in cachedTextEncoder)) {
    cachedTextEncoder.encodeInto = function (arg, view) {
        const buf = cachedTextEncoder.encode(arg);
        view.set(buf);
        return {
            read: arg.length,
            written: buf.length
        };
    }
}

let WASM_VECTOR_LEN = 0;

const AesEncryptedDataFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_aesencrypteddata_free(ptr >>> 0, 1));

const AesEncryptedDataBytesFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_aesencrypteddatabytes_free(ptr >>> 0, 1));

const CommitmentFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_commitment_free(ptr >>> 0, 1));

const DidFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_did_free(ptr >>> 0, 1));

const EncryptedDataFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_encrypteddata_free(ptr >>> 0, 1));

const EncryptedDataBytesFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_encrypteddatabytes_free(ptr >>> 0, 1));

const KeyPairFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_keypair_free(ptr >>> 0, 1));

const TxFinalization = (typeof FinalizationRegistry === 'undefined')
    ? { register: () => {}, unregister: () => {} }
    : new FinalizationRegistry(ptr => wasm.__wbg_tx_free(ptr >>> 0, 1));

export class AesEncryptedData {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(AesEncryptedData.prototype);
        obj.__wbg_ptr = ptr;
        AesEncryptedDataFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        AesEncryptedDataFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_aesencrypteddata_free(ptr, 0);
    }
    /**
     * @param {any} object
     * @returns {AesEncryptedData}
     */
    static from_object(object) {
        const ret = wasm.aesencrypteddata_from_object(object);
        return AesEncryptedData.__wrap(ret);
    }
    /**
     * @param {string} data
     * @param {string} nonce
     */
    constructor(data, nonce) {
        const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(nonce, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.aesencrypteddata_new(ptr0, len0, ptr1, len1);
        this.__wbg_ptr = ret >>> 0;
        AesEncryptedDataFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {string}
     */
    get data() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.aesencrypteddata_data(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get nonce() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.aesencrypteddata_nonce(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.aesencrypteddata_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) AesEncryptedData.prototype[Symbol.dispose] = AesEncryptedData.prototype.free;

export class AesEncryptedDataBytes {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(AesEncryptedDataBytes.prototype);
        obj.__wbg_ptr = ptr;
        AesEncryptedDataBytesFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        AesEncryptedDataBytesFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_aesencrypteddatabytes_free(ptr, 0);
    }
    /**
     * @param {any} object
     * @returns {AesEncryptedDataBytes}
     */
    static from_object(object) {
        const ret = wasm.aesencrypteddatabytes_from_object(object);
        return AesEncryptedDataBytes.__wrap(ret);
    }
    /**
     * @param {Uint8Array} data
     * @param {Uint8Array} nonce
     */
    constructor(data, nonce) {
        const ptr0 = passArray8ToWasm0(data, wasm.__wbindgen_malloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passArray8ToWasm0(nonce, wasm.__wbindgen_malloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.aesencrypteddata_new(ptr0, len0, ptr1, len1);
        this.__wbg_ptr = ret >>> 0;
        AesEncryptedDataBytesFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {Uint8Array}
     */
    get data() {
        const ret = wasm.aesencrypteddatabytes_data(this.__wbg_ptr);
        var v1 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
        wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
        return v1;
    }
    /**
     * @returns {Uint8Array}
     */
    get nonce() {
        const ret = wasm.aesencrypteddatabytes_nonce(this.__wbg_ptr);
        var v1 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
        wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
        return v1;
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.aesencrypteddatabytes_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) AesEncryptedDataBytes.prototype[Symbol.dispose] = AesEncryptedDataBytes.prototype.free;

export class Commitment {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(Commitment.prototype);
        obj.__wbg_ptr = ptr;
        CommitmentFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        CommitmentFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_commitment_free(ptr, 0);
    }
    /**
     * @returns {string}
     */
    get commitment() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.commitment_commitment(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {any} object
     * @returns {Commitment}
     */
    static from_object(object) {
        const ret = wasm.commitment_from_object(object);
        return Commitment.__wrap(ret);
    }
    /**
     * @param {string} commitment
     */
    set commitment(commitment) {
        const ptr0 = passStringToWasm0(commitment, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.commitment_set_commitment(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @returns {string}
     */
    get secret_blinding_factor() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.commitment_secret_blinding_factor(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {string} secret_blinding_factor
     */
    set secret_blinding_factor(secret_blinding_factor) {
        const ptr0 = passStringToWasm0(secret_blinding_factor, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.commitment_set_secret_blinding_factor(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {string} commitment
     * @param {string} secret_blinding_factor
     */
    constructor(commitment, secret_blinding_factor) {
        const ptr0 = passStringToWasm0(commitment, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(secret_blinding_factor, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.aesencrypteddata_new(ptr0, len0, ptr1, len1);
        this.__wbg_ptr = ret >>> 0;
        CommitmentFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.commitment_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) Commitment.prototype[Symbol.dispose] = Commitment.prototype.free;

export class Did {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(Did.prototype);
        obj.__wbg_ptr = ptr;
        DidFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        DidFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_did_free(ptr, 0);
    }
    /**
     * @param {any} object
     * @returns {Did}
     */
    static from_object(object) {
        const ret = wasm.did_from_object(object);
        return Did.__wrap(ret);
    }
    /**
     * @returns {string}
     */
    get id() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.did_id(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {string} id
     */
    constructor(id) {
        const ptr0 = passStringToWasm0(id, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ret = wasm.did_new(ptr0, len0);
        this.__wbg_ptr = ret >>> 0;
        DidFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @param {string} id
     */
    set id(id) {
        const ptr0 = passStringToWasm0(id, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.did_set_id(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.did_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) Did.prototype[Symbol.dispose] = Did.prototype.free;

/**
 * @enum {0 | 1 | 2}
 */
export const EncodingType = Object.freeze({
    UTF8: 0, "0": "UTF8",
    HEX: 1, "1": "HEX",
    BASE64: 2, "2": "BASE64",
});

export class EncryptedData {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(EncryptedData.prototype);
        obj.__wbg_ptr = ptr;
        EncryptedDataFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        EncryptedDataFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_encrypteddata_free(ptr, 0);
    }
    /**
     * @param {any} object
     * @returns {EncryptedData}
     */
    static from_object(object) {
        const ret = wasm.encrypteddata_from_object(object);
        return EncryptedData.__wrap(ret);
    }
    /**
     * @returns {string}
     */
    get pubkey_from() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.encrypteddata_pubkey_from(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {string} pubkey_from
     * @param {string} pubkey_to
     * @param {string} data
     * @param {string} nonce
     */
    constructor(pubkey_from, pubkey_to, data, nonce) {
        const ptr0 = passStringToWasm0(pubkey_from, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(pubkey_to, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ptr2 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len2 = WASM_VECTOR_LEN;
        const ptr3 = passStringToWasm0(nonce, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len3 = WASM_VECTOR_LEN;
        const ret = wasm.encrypteddata_new(ptr0, len0, ptr1, len1, ptr2, len2, ptr3, len3);
        this.__wbg_ptr = ret >>> 0;
        EncryptedDataFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {string}
     */
    get data() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.encrypteddata_data(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get nonce() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.encrypteddata_nonce(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get pubkey_to() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.encrypteddata_pubkey_to(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.encrypteddata_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) EncryptedData.prototype[Symbol.dispose] = EncryptedData.prototype.free;

export class EncryptedDataBytes {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(EncryptedDataBytes.prototype);
        obj.__wbg_ptr = ptr;
        EncryptedDataBytesFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        EncryptedDataBytesFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_encrypteddatabytes_free(ptr, 0);
    }
    /**
     * @param {any} object
     * @returns {EncryptedDataBytes}
     */
    static from_object(object) {
        const ret = wasm.encrypteddatabytes_from_object(object);
        return EncryptedDataBytes.__wrap(ret);
    }
    /**
     * @returns {string}
     */
    get pubkey_from() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.encrypteddatabytes_pubkey_from(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {string} pubkey_from
     * @param {string} pubkey_to
     * @param {Uint8Array} data
     * @param {Uint8Array} nonce
     */
    constructor(pubkey_from, pubkey_to, data, nonce) {
        const ptr0 = passStringToWasm0(pubkey_from, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(pubkey_to, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ptr2 = passArray8ToWasm0(data, wasm.__wbindgen_malloc);
        const len2 = WASM_VECTOR_LEN;
        const ptr3 = passArray8ToWasm0(nonce, wasm.__wbindgen_malloc);
        const len3 = WASM_VECTOR_LEN;
        const ret = wasm.encrypteddata_new(ptr0, len0, ptr1, len1, ptr2, len2, ptr3, len3);
        this.__wbg_ptr = ret >>> 0;
        EncryptedDataBytesFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {Uint8Array}
     */
    get data() {
        const ret = wasm.encrypteddatabytes_data(this.__wbg_ptr);
        var v1 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
        wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
        return v1;
    }
    /**
     * @returns {Uint8Array}
     */
    get nonce() {
        const ret = wasm.encrypteddatabytes_nonce(this.__wbg_ptr);
        var v1 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
        wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
        return v1;
    }
    /**
     * @returns {string}
     */
    get pubkey_to() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.encrypteddatabytes_pubkey_to(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.encrypteddatabytes_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) EncryptedDataBytes.prototype[Symbol.dispose] = EncryptedDataBytes.prototype.free;

export class KeyPair {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(KeyPair.prototype);
        obj.__wbg_ptr = ptr;
        KeyPairFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        KeyPairFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_keypair_free(ptr, 0);
    }
    /**
     * @param {string} pubkey
     */
    set pubkey(pubkey) {
        const ptr0 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.keypair_set_pubkey(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {any} object
     * @returns {KeyPair}
     */
    static from_object(object) {
        const ret = wasm.keypair_from_object(object);
        return KeyPair.__wrap(ret);
    }
    /**
     * @param {string} privkey
     */
    set privkey(privkey) {
        const ptr0 = passStringToWasm0(privkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.keypair_set_privkey(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {string} pubkey
     * @param {string} privkey
     */
    constructor(pubkey, privkey) {
        const ptr0 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(privkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.aesencrypteddata_new(ptr0, len0, ptr1, len1);
        this.__wbg_ptr = ret >>> 0;
        KeyPairFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {string}
     */
    get pubkey() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.keypair_pubkey(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get privkey() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.keypair_privkey(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.keypair_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) KeyPair.prototype[Symbol.dispose] = KeyPair.prototype.free;

export class Tx {
    static __wrap(ptr) {
        ptr = ptr >>> 0;
        const obj = Object.create(Tx.prototype);
        obj.__wbg_ptr = ptr;
        TxFinalization.register(obj, obj.__wbg_ptr, obj);
        return obj;
    }
    __destroy_into_raw() {
        const ptr = this.__wbg_ptr;
        this.__wbg_ptr = 0;
        TxFinalization.unregister(this);
        return ptr;
    }
    free() {
        const ptr = this.__destroy_into_raw();
        wasm.__wbg_tx_free(ptr, 0);
    }
    /**
     * @param {string} amount
     */
    set amount(amount) {
        const ptr0 = passStringToWasm0(amount, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.tx_set_amount(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {any} object
     * @returns {Tx}
     */
    static from_object(object) {
        const ret = wasm.tx_from_object(object);
        return Tx.__wrap(ret);
    }
    /**
     * @returns {string}
     */
    get to() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.tx_to(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get fee() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.tx_fee(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {string} coin
     * @param {string} from
     * @param {string} to
     * @param {string} amount
     * @param {string} fee
     * @param {string} data
     */
    constructor(coin, from, to, amount, fee, data) {
        const ptr0 = passStringToWasm0(coin, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(from, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ptr2 = passStringToWasm0(to, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len2 = WASM_VECTOR_LEN;
        const ptr3 = passStringToWasm0(amount, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len3 = WASM_VECTOR_LEN;
        const ptr4 = passStringToWasm0(fee, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len4 = WASM_VECTOR_LEN;
        const ptr5 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len5 = WASM_VECTOR_LEN;
        const ret = wasm.tx_new(ptr0, len0, ptr1, len1, ptr2, len2, ptr3, len3, ptr4, len4, ptr5, len5);
        this.__wbg_ptr = ret >>> 0;
        TxFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
    /**
     * @returns {string}
     */
    get coin() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.tx_coin(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get data() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.tx_data(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get from() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.tx_from(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @returns {string}
     */
    get amount() {
        let deferred1_0;
        let deferred1_1;
        try {
            const ret = wasm.tx_amount(this.__wbg_ptr);
            deferred1_0 = ret[0];
            deferred1_1 = ret[1];
            return getStringFromWasm0(ret[0], ret[1]);
        } finally {
            wasm.__wbindgen_free(deferred1_0, deferred1_1, 1);
        }
    }
    /**
     * @param {string} to
     */
    set to(to) {
        const ptr0 = passStringToWasm0(to, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.tx_set_to(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {string} fee
     */
    set fee(fee) {
        const ptr0 = passStringToWasm0(fee, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.tx_set_fee(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {string} coin
     */
    set coin(coin) {
        const ptr0 = passStringToWasm0(coin, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.tx_set_coin(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {string} data
     */
    set data(data) {
        const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.tx_set_data(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @param {string} from
     */
    set from(from) {
        const ptr0 = passStringToWasm0(from, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        wasm.tx_set_from(this.__wbg_ptr, ptr0, len0);
    }
    /**
     * @returns {any}
     */
    to_object() {
        const ret = wasm.tx_to_object(this.__wbg_ptr);
        return ret;
    }
}
if (Symbol.dispose) Tx.prototype[Symbol.dispose] = Tx.prototype.free;

/**
 * @returns {KeyPair}
 */
export function create_keypair() {
    const ret = wasm.create_keypair();
    return KeyPair.__wrap(ret);
}

/**
 * @param {EncryptedData} data
 * @param {string} privkey
 * @param {EncodingType} encoding_type
 * @returns {string}
 */
export function decrypt(data, privkey, encoding_type) {
    let deferred4_0;
    let deferred4_1;
    try {
        _assertClass(data, EncryptedData);
        var ptr0 = data.__destroy_into_raw();
        const ptr1 = passStringToWasm0(privkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.decrypt(ptr0, ptr1, len1, encoding_type);
        var ptr3 = ret[0];
        var len3 = ret[1];
        if (ret[3]) {
            ptr3 = 0; len3 = 0;
            throw takeFromExternrefTable0(ret[2]);
        }
        deferred4_0 = ptr3;
        deferred4_1 = len3;
        return getStringFromWasm0(ptr3, len3);
    } finally {
        wasm.__wbindgen_free(deferred4_0, deferred4_1, 1);
    }
}

/**
 * @param {AesEncryptedData} data
 * @param {string} key
 * @param {EncodingType} encoding_type
 * @returns {string}
 */
export function decrypt_aes(data, key, encoding_type) {
    let deferred4_0;
    let deferred4_1;
    try {
        _assertClass(data, AesEncryptedData);
        var ptr0 = data.__destroy_into_raw();
        const ptr1 = passStringToWasm0(key, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.decrypt_aes(ptr0, ptr1, len1, encoding_type);
        var ptr3 = ret[0];
        var len3 = ret[1];
        if (ret[3]) {
            ptr3 = 0; len3 = 0;
            throw takeFromExternrefTable0(ret[2]);
        }
        deferred4_0 = ptr3;
        deferred4_1 = len3;
        return getStringFromWasm0(ptr3, len3);
    } finally {
        wasm.__wbindgen_free(deferred4_0, deferred4_1, 1);
    }
}

/**
 * @param {AesEncryptedDataBytes} data
 * @param {string} key
 * @returns {Uint8Array}
 */
export function decrypt_aes_bytes(data, key) {
    _assertClass(data, AesEncryptedDataBytes);
    var ptr0 = data.__destroy_into_raw();
    const ptr1 = passStringToWasm0(key, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.decrypt_aes_bytes(ptr0, ptr1, len1);
    if (ret[3]) {
        throw takeFromExternrefTable0(ret[2]);
    }
    var v3 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
    wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
    return v3;
}

/**
 * @param {EncryptedDataBytes} data
 * @param {string} privkey
 * @returns {Uint8Array}
 */
export function decrypt_bytes(data, privkey) {
    _assertClass(data, EncryptedDataBytes);
    var ptr0 = data.__destroy_into_raw();
    const ptr1 = passStringToWasm0(privkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.decrypt_bytes(ptr0, ptr1, len1);
    if (ret[3]) {
        throw takeFromExternrefTable0(ret[2]);
    }
    var v3 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
    wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
    return v3;
}

/**
 * @param {Did} did
 * @returns {string}
 */
export function did_to_key(did) {
    let deferred3_0;
    let deferred3_1;
    try {
        _assertClass(did, Did);
        var ptr0 = did.__destroy_into_raw();
        const ret = wasm.did_to_key(ptr0);
        var ptr2 = ret[0];
        var len2 = ret[1];
        if (ret[3]) {
            ptr2 = 0; len2 = 0;
            throw takeFromExternrefTable0(ret[2]);
        }
        deferred3_0 = ptr2;
        deferred3_1 = len2;
        return getStringFromWasm0(ptr2, len2);
    } finally {
        wasm.__wbindgen_free(deferred3_0, deferred3_1, 1);
    }
}

/**
 * @param {string} privkey
 * @param {string} pubkey
 * @returns {string}
 */
export function ecdh(privkey, pubkey) {
    let deferred4_0;
    let deferred4_1;
    try {
        const ptr0 = passStringToWasm0(privkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.ecdh(ptr0, len0, ptr1, len1);
        var ptr3 = ret[0];
        var len3 = ret[1];
        if (ret[3]) {
            ptr3 = 0; len3 = 0;
            throw takeFromExternrefTable0(ret[2]);
        }
        deferred4_0 = ptr3;
        deferred4_1 = len3;
        return getStringFromWasm0(ptr3, len3);
    } finally {
        wasm.__wbindgen_free(deferred4_0, deferred4_1, 1);
    }
}

/**
 * @param {string} data
 * @param {string} pubkey
 * @param {EncodingType} encoding_type
 * @returns {EncryptedData}
 */
export function encrypt(data, pubkey, encoding_type) {
    const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len0 = WASM_VECTOR_LEN;
    const ptr1 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.encrypt(ptr0, len0, ptr1, len1, encoding_type);
    if (ret[2]) {
        throw takeFromExternrefTable0(ret[1]);
    }
    return EncryptedData.__wrap(ret[0]);
}

/**
 * @param {string} data
 * @param {string} key
 * @param {EncodingType} encoding_type
 * @returns {AesEncryptedData}
 */
export function encrypt_aes(data, key, encoding_type) {
    const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len0 = WASM_VECTOR_LEN;
    const ptr1 = passStringToWasm0(key, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.encrypt_aes(ptr0, len0, ptr1, len1, encoding_type);
    if (ret[2]) {
        throw takeFromExternrefTable0(ret[1]);
    }
    return AesEncryptedData.__wrap(ret[0]);
}

/**
 * @param {Uint8Array} data
 * @param {string} key
 * @returns {AesEncryptedDataBytes}
 */
export function encrypt_aes_bytes(data, key) {
    const ptr0 = passArray8ToWasm0(data, wasm.__wbindgen_malloc);
    const len0 = WASM_VECTOR_LEN;
    const ptr1 = passStringToWasm0(key, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.encrypt_aes_bytes(ptr0, len0, ptr1, len1);
    if (ret[2]) {
        throw takeFromExternrefTable0(ret[1]);
    }
    return AesEncryptedDataBytes.__wrap(ret[0]);
}

/**
 * @param {Uint8Array} data
 * @param {string} pubkey
 * @returns {EncryptedDataBytes}
 */
export function encrypt_bytes(data, pubkey) {
    const ptr0 = passArray8ToWasm0(data, wasm.__wbindgen_malloc);
    const len0 = WASM_VECTOR_LEN;
    const ptr1 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.encrypt_bytes(ptr0, len0, ptr1, len1);
    if (ret[2]) {
        throw takeFromExternrefTable0(ret[1]);
    }
    return EncryptedDataBytes.__wrap(ret[0]);
}

export function init_panic_hook() {
    wasm.init_panic_hook();
}

/**
 * @param {string} pubkey
 * @returns {Did}
 */
export function key_to_did(pubkey) {
    const ptr0 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len0 = WASM_VECTOR_LEN;
    const ret = wasm.key_to_did(ptr0, len0);
    return Did.__wrap(ret);
}

/**
 * @param {bigint} value
 * @param {string} tag
 * @returns {Commitment}
 */
export function pedersen_commit(value, tag) {
    const ptr0 = passStringToWasm0(tag, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len0 = WASM_VECTOR_LEN;
    const ret = wasm.pedersen_commit(value, ptr0, len0);
    return Commitment.__wrap(ret);
}

/**
 * @param {Commitment} commitment
 * @param {bigint} value
 * @param {string} tag
 * @returns {boolean}
 */
export function pedersen_reveal(commitment, value, tag) {
    _assertClass(commitment, Commitment);
    var ptr0 = commitment.__destroy_into_raw();
    const ptr1 = passStringToWasm0(tag, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ret = wasm.pedersen_reveal(ptr0, value, ptr1, len1);
    if (ret[2]) {
        throw takeFromExternrefTable0(ret[1]);
    }
    return ret[0] !== 0;
}

/**
 * @param {string} data
 * @returns {string}
 */
export function sha256(data) {
    let deferred2_0;
    let deferred2_1;
    try {
        const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ret = wasm.sha256(ptr0, len0);
        deferred2_0 = ret[0];
        deferred2_1 = ret[1];
        return getStringFromWasm0(ret[0], ret[1]);
    } finally {
        wasm.__wbindgen_free(deferred2_0, deferred2_1, 1);
    }
}

/**
 * @param {Uint8Array} data
 * @returns {Uint8Array}
 */
export function sha256_bytes(data) {
    const ptr0 = passArray8ToWasm0(data, wasm.__wbindgen_malloc);
    const len0 = WASM_VECTOR_LEN;
    const ret = wasm.sha256_bytes(ptr0, len0);
    var v2 = getArrayU8FromWasm0(ret[0], ret[1]).slice();
    wasm.__wbindgen_free(ret[0], ret[1] * 1, 1);
    return v2;
}

/**
 * @param {string} data
 * @param {string} privkey
 * @returns {string}
 */
export function sign(data, privkey) {
    let deferred4_0;
    let deferred4_1;
    try {
        const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(privkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ret = wasm.sign(ptr0, len0, ptr1, len1);
        var ptr3 = ret[0];
        var len3 = ret[1];
        if (ret[3]) {
            ptr3 = 0; len3 = 0;
            throw takeFromExternrefTable0(ret[2]);
        }
        deferred4_0 = ptr3;
        deferred4_1 = len3;
        return getStringFromWasm0(ptr3, len3);
    } finally {
        wasm.__wbindgen_free(deferred4_0, deferred4_1, 1);
    }
}

/**
 * @param {string} data
 * @param {string} sig
 * @param {string} pubkey
 * @returns {boolean}
 */
export function verify(data, sig, pubkey) {
    const ptr0 = passStringToWasm0(data, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len0 = WASM_VECTOR_LEN;
    const ptr1 = passStringToWasm0(sig, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    const ptr2 = passStringToWasm0(pubkey, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len2 = WASM_VECTOR_LEN;
    const ret = wasm.verify(ptr0, len0, ptr1, len1, ptr2, len2);
    if (ret[2]) {
        throw takeFromExternrefTable0(ret[1]);
    }
    return ret[0] !== 0;
}

export function __wbg_Error_52673b7de5a0ca89(arg0, arg1) {
    const ret = Error(getStringFromWasm0(arg0, arg1));
    return ret;
};

export function __wbg_Number_2d1dcfcf4ec51736(arg0) {
    const ret = Number(arg0);
    return ret;
};

export function __wbg_String_8f0eb39a4a4c2f66(arg0, arg1) {
    const ret = String(arg1);
    const ptr1 = passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    getDataViewMemory0().setInt32(arg0 + 4 * 1, len1, true);
    getDataViewMemory0().setInt32(arg0 + 4 * 0, ptr1, true);
};

export function __wbg___wbindgen_boolean_get_dea25b33882b895b(arg0) {
    const v = arg0;
    const ret = typeof(v) === 'boolean' ? v : undefined;
    return isLikeNone(ret) ? 0xFFFFFF : ret ? 1 : 0;
};

export function __wbg___wbindgen_debug_string_adfb662ae34724b6(arg0, arg1) {
    const ret = debugString(arg1);
    const ptr1 = passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    getDataViewMemory0().setInt32(arg0 + 4 * 1, len1, true);
    getDataViewMemory0().setInt32(arg0 + 4 * 0, ptr1, true);
};

export function __wbg___wbindgen_in_0d3e1e8f0c669317(arg0, arg1) {
    const ret = arg0 in arg1;
    return ret;
};

export function __wbg___wbindgen_is_function_8d400b8b1af978cd(arg0) {
    const ret = typeof(arg0) === 'function';
    return ret;
};

export function __wbg___wbindgen_is_object_ce774f3490692386(arg0) {
    const val = arg0;
    const ret = typeof(val) === 'object' && val !== null;
    return ret;
};

export function __wbg___wbindgen_is_string_704ef9c8fc131030(arg0) {
    const ret = typeof(arg0) === 'string';
    return ret;
};

export function __wbg___wbindgen_is_undefined_f6b95eab589e0269(arg0) {
    const ret = arg0 === undefined;
    return ret;
};

export function __wbg___wbindgen_jsval_loose_eq_766057600fdd1b0d(arg0, arg1) {
    const ret = arg0 == arg1;
    return ret;
};

export function __wbg___wbindgen_number_get_9619185a74197f95(arg0, arg1) {
    const obj = arg1;
    const ret = typeof(obj) === 'number' ? obj : undefined;
    getDataViewMemory0().setFloat64(arg0 + 8 * 1, isLikeNone(ret) ? 0 : ret, true);
    getDataViewMemory0().setInt32(arg0 + 4 * 0, !isLikeNone(ret), true);
};

export function __wbg___wbindgen_string_get_a2a31e16edf96e42(arg0, arg1) {
    const obj = arg1;
    const ret = typeof(obj) === 'string' ? obj : undefined;
    var ptr1 = isLikeNone(ret) ? 0 : passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    var len1 = WASM_VECTOR_LEN;
    getDataViewMemory0().setInt32(arg0 + 4 * 1, len1, true);
    getDataViewMemory0().setInt32(arg0 + 4 * 0, ptr1, true);
};

export function __wbg___wbindgen_throw_dd24417ed36fc46e(arg0, arg1) {
    throw new Error(getStringFromWasm0(arg0, arg1));
};

export function __wbg_call_3020136f7a2d6e44() { return handleError(function (arg0, arg1, arg2) {
    const ret = arg0.call(arg1, arg2);
    return ret;
}, arguments) };

export function __wbg_call_abb4ff46ce38be40() { return handleError(function (arg0, arg1) {
    const ret = arg0.call(arg1);
    return ret;
}, arguments) };

export function __wbg_crypto_ed58b8e10a292839(arg0) {
    const ret = arg0.crypto;
    return ret;
};

export function __wbg_done_62ea16af4ce34b24(arg0) {
    const ret = arg0.done;
    return ret;
};

export function __wbg_error_7534b8e9a36f1ab4(arg0, arg1) {
    let deferred0_0;
    let deferred0_1;
    try {
        deferred0_0 = arg0;
        deferred0_1 = arg1;
        console.error(getStringFromWasm0(arg0, arg1));
    } finally {
        wasm.__wbindgen_free(deferred0_0, deferred0_1, 1);
    }
};

export function __wbg_getRandomValues_bcb4912f16000dc4() { return handleError(function (arg0, arg1) {
    arg0.getRandomValues(arg1);
}, arguments) };

export function __wbg_get_6b7bd52aca3f9671(arg0, arg1) {
    const ret = arg0[arg1 >>> 0];
    return ret;
};

export function __wbg_get_af9dab7e9603ea93() { return handleError(function (arg0, arg1) {
    const ret = Reflect.get(arg0, arg1);
    return ret;
}, arguments) };

export function __wbg_get_with_ref_key_1dc361bd10053bfe(arg0, arg1) {
    const ret = arg0[arg1];
    return ret;
};

export function __wbg_instanceof_ArrayBuffer_f3320d2419cd0355(arg0) {
    let result;
    try {
        result = arg0 instanceof ArrayBuffer;
    } catch (_) {
        result = false;
    }
    const ret = result;
    return ret;
};

export function __wbg_instanceof_Uint8Array_da54ccc9d3e09434(arg0) {
    let result;
    try {
        result = arg0 instanceof Uint8Array;
    } catch (_) {
        result = false;
    }
    const ret = result;
    return ret;
};

export function __wbg_isArray_51fd9e6422c0a395(arg0) {
    const ret = Array.isArray(arg0);
    return ret;
};

export function __wbg_isSafeInteger_ae7d3f054d55fa16(arg0) {
    const ret = Number.isSafeInteger(arg0);
    return ret;
};

export function __wbg_iterator_27b7c8b35ab3e86b() {
    const ret = Symbol.iterator;
    return ret;
};

export function __wbg_length_22ac23eaec9d8053(arg0) {
    const ret = arg0.length;
    return ret;
};

export function __wbg_length_d45040a40c570362(arg0) {
    const ret = arg0.length;
    return ret;
};

export function __wbg_msCrypto_0a36e2ec3a343d26(arg0) {
    const ret = arg0.msCrypto;
    return ret;
};

export function __wbg_new_1ba21ce319a06297() {
    const ret = new Object();
    return ret;
};

export function __wbg_new_25f239778d6112b9() {
    const ret = new Array();
    return ret;
};

export function __wbg_new_6421f6084cc5bc5a(arg0) {
    const ret = new Uint8Array(arg0);
    return ret;
};

export function __wbg_new_8a6f238a6ece86ea() {
    const ret = new Error();
    return ret;
};

export function __wbg_new_no_args_cb138f77cf6151ee(arg0, arg1) {
    const ret = new Function(getStringFromWasm0(arg0, arg1));
    return ret;
};

export function __wbg_new_with_length_aa5eaf41d35235e5(arg0) {
    const ret = new Uint8Array(arg0 >>> 0);
    return ret;
};

export function __wbg_next_138a17bbf04e926c(arg0) {
    const ret = arg0.next;
    return ret;
};

export function __wbg_next_3cfe5c0fe2a4cc53() { return handleError(function (arg0) {
    const ret = arg0.next();
    return ret;
}, arguments) };

export function __wbg_node_02999533c4ea02e3(arg0) {
    const ret = arg0.node;
    return ret;
};

export function __wbg_process_5c1d670bc53614b8(arg0) {
    const ret = arg0.process;
    return ret;
};

export function __wbg_prototypesetcall_dfe9b766cdc1f1fd(arg0, arg1, arg2) {
    Uint8Array.prototype.set.call(getArrayU8FromWasm0(arg0, arg1), arg2);
};

export function __wbg_randomFillSync_ab2cfe79ebbf2740() { return handleError(function (arg0, arg1) {
    arg0.randomFillSync(arg1);
}, arguments) };

export function __wbg_require_79b1e9274cde3c87() { return handleError(function () {
    const ret = module.require;
    return ret;
}, arguments) };

export function __wbg_set_3f1d0b984ed272ed(arg0, arg1, arg2) {
    arg0[arg1] = arg2;
};

export function __wbg_set_7df433eea03a5c14(arg0, arg1, arg2) {
    arg0[arg1 >>> 0] = arg2;
};

export function __wbg_stack_0ed75d68575b0f3c(arg0, arg1) {
    const ret = arg1.stack;
    const ptr1 = passStringToWasm0(ret, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
    const len1 = WASM_VECTOR_LEN;
    getDataViewMemory0().setInt32(arg0 + 4 * 1, len1, true);
    getDataViewMemory0().setInt32(arg0 + 4 * 0, ptr1, true);
};

export function __wbg_static_accessor_GLOBAL_769e6b65d6557335() {
    const ret = typeof global === 'undefined' ? null : global;
    return isLikeNone(ret) ? 0 : addToExternrefTable0(ret);
};

export function __wbg_static_accessor_GLOBAL_THIS_60cf02db4de8e1c1() {
    const ret = typeof globalThis === 'undefined' ? null : globalThis;
    return isLikeNone(ret) ? 0 : addToExternrefTable0(ret);
};

export function __wbg_static_accessor_SELF_08f5a74c69739274() {
    const ret = typeof self === 'undefined' ? null : self;
    return isLikeNone(ret) ? 0 : addToExternrefTable0(ret);
};

export function __wbg_static_accessor_WINDOW_a8924b26aa92d024() {
    const ret = typeof window === 'undefined' ? null : window;
    return isLikeNone(ret) ? 0 : addToExternrefTable0(ret);
};

export function __wbg_subarray_845f2f5bce7d061a(arg0, arg1, arg2) {
    const ret = arg0.subarray(arg1 >>> 0, arg2 >>> 0);
    return ret;
};

export function __wbg_value_57b7b035e117f7ee(arg0) {
    const ret = arg0.value;
    return ret;
};

export function __wbg_versions_c71aa1626a93e0a1(arg0) {
    const ret = arg0.versions;
    return ret;
};

export function __wbindgen_cast_2241b6af4c4b2941(arg0, arg1) {
    // Cast intrinsic for `Ref(String) -> Externref`.
    const ret = getStringFromWasm0(arg0, arg1);
    return ret;
};

export function __wbindgen_cast_cb9088102bce6b30(arg0, arg1) {
    // Cast intrinsic for `Ref(Slice(U8)) -> NamedExternref("Uint8Array")`.
    const ret = getArrayU8FromWasm0(arg0, arg1);
    return ret;
};

export function __wbindgen_cast_d6cd19b81560fd6e(arg0) {
    // Cast intrinsic for `F64 -> Externref`.
    const ret = arg0;
    return ret;
};

export function __wbindgen_init_externref_table() {
    const table = wasm.__wbindgen_externrefs;
    const offset = table.grow(4);
    table.set(0, undefined);
    table.set(offset + 0, undefined);
    table.set(offset + 1, null);
    table.set(offset + 2, true);
    table.set(offset + 3, false);
};
