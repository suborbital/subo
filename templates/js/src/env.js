export const LogLevel = Object.freeze({
    0: "Null",
    "Null": 0,
    1: "Error",
    "Error": 1,
    2: "Warn",
    "Warn": 2,
    3: "Info",
    "Info": 3,
    4: "Debug",
    "Debug": 4,
  });
  export const HttpMethod = Object.freeze({
    0: "Get",
    "Get": 0,
    1: "Head",
    "Head": 1,
    2: "Options",
    "Options": 2,
    3: "Post",
    "Post": 3,
    4: "Put",
    "Put": 4,
    5: "Patch",
    "Patch": 5,
    6: "Delete",
    "Delete": 6,
  });
  export const FieldType = Object.freeze({
    0: "Meta",
    "Meta": 0,
    1: "Body",
    "Body": 1,
    2: "Header",
    "Header": 2,
    3: "Params",
    "Params": 3,
    4: "State",
    "State": 4,
    5: "Query",
    "Query": 5,
  });
  function clamp_host(i, min, max) {
    if (!Number.isInteger(i)) throw new TypeError(`must be an integer`);
    if (i < min || i > max) throw new RangeError(`must be between ${min} and ${max}`);
    return i;
  }
  
  let UTF8_ENCODED_LEN = 0;
  const UTF8_ENCODER = new TextEncoder('utf-8');
  
  function utf8_encode(s, realloc, memory) {
    if (typeof s !== 'string') throw new TypeError('expected a string');
    
    if (s.length === 0) {
      UTF8_ENCODED_LEN = 0;
      return 1;
    }
    
    let alloc_len = 0;
    let ptr = 0;
    let writtenTotal = 0;
    while (s.length > 0) {
      ptr = realloc(ptr, alloc_len, 1, alloc_len + s.length);
        alloc_len += s.length;
      let buf = new Uint8Array(memory.buffer, ptr + writtenTotal, alloc_len - writtenTotal)
      const { read, written } = UTF8_ENCODER.encodeInto(
      s,
      buf,
      );
      writtenTotal += written;
      s = s.slice(read);
    }
    if (alloc_len > writtenTotal)
    ptr = realloc(ptr, alloc_len, 1, writtenTotal);
    UTF8_ENCODED_LEN = writtenTotal;
    return ptr;
  }
  export class Env {
    addToImports(imports) {
    }
    
    async instantiate(module, imports) {
      imports = imports || {};
      this.addToImports(imports);
      
      if (module instanceof WebAssembly.Instance) {
        this.instance = module;
      } else if (module instanceof WebAssembly.Module) {
        this.instance = await WebAssembly.instantiate(module, imports);
      } else if (module instanceof ArrayBuffer || module instanceof Uint8Array) {
        const { instance } = await WebAssembly.instantiate(module, imports);
        this.instance = instance;
      } else {
        const { instance } = await WebAssembly.instantiateStreaming(module, imports);
        this.instance = instance;
      }
      this._exports = this.instance.exports;
    }
    returnResult(arg0, arg1) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const val0 = arg0;
      const len0 = val0.length;
      const ptr0 = realloc(0, 0, 1, len0 * 1);
      (new Uint8Array(memory.buffer, ptr0, len0 * 1)).set(new Uint8Array(val0.buffer));
      this._exports['return_result'](ptr0, len0, clamp_host(arg1, 0, 4294967295));
    }
    returnError(arg0, arg1, arg2) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const ptr0 = utf8_encode(arg1, realloc, memory);
      const len0 = UTF8_ENCODED_LEN;
      this._exports['return_error'](clamp_host(arg0, 0, 4294967295), ptr0, len0, clamp_host(arg2, 0, 4294967295));
    }
    logMsg(arg0, arg1, arg2) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const ptr0 = utf8_encode.call(this, arg0, realloc, memory);
      const len0 = UTF8_ENCODED_LEN;
      const variant1 = arg1;
      if (!(variant1 in LogLevel))
      throw new RangeError("invalid variant specified for LogLevel");
      this._exports['log_msg'](ptr0, len0, Number.isInteger(variant1) ? variant1 : LogLevel[variant1], clamp_host(arg2, 0, 4294967295));
    }
    fetchUrl(arg0, arg1, arg2, arg3) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const variant0 = arg0;
      if (!(variant0 in HttpMethod))
      throw new RangeError("invalid variant specified for HttpMethod");
      const ptr1 = utf8_encode(arg1, realloc, memory);
      const len1 = UTF8_ENCODED_LEN;
      const val2 = arg2;
      const len2 = val2.length;
      const ptr2 = realloc(0, 0, 1, len2 * 1);
      (new Uint8Array(memory.buffer, ptr2, len2 * 1)).set(new Uint8Array(val2.buffer));
      const ret = this._exports['fetch_url'](Number.isInteger(variant0) ? variant0 : HttpMethod[variant0], ptr1, len1, ptr2, len2, clamp_host(arg3, 0, 4294967295));
      return ret >>> 0;
    }
    graphqlQuery(arg0, arg1, arg2) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const ptr0 = utf8_encode(arg0, realloc, memory);
      const len0 = UTF8_ENCODED_LEN;
      const ptr1 = utf8_encode(arg1, realloc, memory);
      const len1 = UTF8_ENCODED_LEN;
      const ret = this._exports['graphql_query'](ptr0, len0, ptr1, len1, clamp_host(arg2, 0, 4294967295));
      return ret >>> 0;
    }
    cacheSet(arg0, arg1, arg2, arg3) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const ptr0 = utf8_encode(arg0, realloc, memory);
      const len0 = UTF8_ENCODED_LEN;
      const val1 = arg1;
      const len1 = val1.length;
      const ptr1 = realloc(0, 0, 1, len1 * 1);
      (new Uint8Array(memory.buffer, ptr1, len1 * 1)).set(new Uint8Array(val1.buffer));
      const ret = this._exports['cache_set'](ptr0, len0, ptr1, len1, clamp_host(arg2, 0, 4294967295), clamp_host(arg3, 0, 4294967295));
      return ret >>> 0;
    }
    cacheGet(arg0, arg1) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const ptr0 = utf8_encode(arg0, realloc, memory);
      const len0 = UTF8_ENCODED_LEN;
      const ret = this._exports['cache_get'](ptr0, len0, clamp_host(arg1, 0, 4294967295));
      return ret >>> 0;
    }
    requestGetField(arg0, arg1, arg2) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const variant0 = arg0;
      if (!(variant0 in FieldType))
      throw new RangeError("invalid variant specified for FieldType");
      const ptr1 = utf8_encode(arg1, realloc, memory);
      const len1 = UTF8_ENCODED_LEN;
      const ret = this._exports['request_get_field'](Number.isInteger(variant0) ? variant0 : FieldType[variant0], ptr1, len1, clamp_host(arg2, 0, 4294967295));
      return ret >>> 0;
    }
    getFfiResult(arg0, arg1) {
      const ret = this._exports['get_ffi_result'](clamp_host(arg0, 0, 4294967295), clamp_host(arg1, 0, 4294967295));
      return ret >>> 0;
    }
    returnAbort(arg0, arg1, arg2, arg3, arg4) {
      const memory = this._exports.memory;
      const realloc = this._exports["canonical_abi_realloc"];
      const ptr0 = utf8_encode(arg0, realloc, memory);
      const len0 = UTF8_ENCODED_LEN;
      const ptr1 = utf8_encode(arg1, realloc, memory);
      const len1 = UTF8_ENCODED_LEN;
      this._exports['return_abort'](ptr0, len0, ptr1, len1, clamp_host(arg2, 0, 4294967295), clamp_host(arg3, 0, 4294967295), clamp_host(arg4, 0, 4294967295));
    }
  }
  