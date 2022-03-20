const RequestFieldTypeMeta   = 0
const RequestFieldTypeBody   = 1
const RequestFieldTypeHeader = 2
const RequestFieldTypeParams = 3
const RequestFieldTypeState  = 4

const textDecoder = new TextDecoder;
// FIXME: switch to KV cache or DO
const cache = {};

function generateRandomId() {
  return Math.floor(Math.random() * Number.MAX_SAFE_INTEGER);
}

async function runRunnable(jobBytes, params, method) {
  let exports;
  let response = new Response("runnable didn't return a response", { status: 404 });
  let ffiResult = null;
  let jobParams = {};

  const ident = generateRandomId();

  function allocate(size) {
    const ptr = exports.allocate(size);
    if (ptr === 0) {
      throw new Error("failed to allocate");
    }
    return ptr;
  }

  function writeData(ptr, data) {
    const view = new Uint8Array(exports.memory.buffer);
    for (let i = 0, len = data.length; i < len; i++) {
      view[ptr + i] = data.charCodeAt(i);
    }
  }

  function readData(ptr, len) {
    const view = new Uint8Array(exports.memory.buffer);
    return view.slice(ptr, ptr+len);
  }

  function readString(ptr, len) {
    const data = readData(ptr, len);
    return textDecoder.decode(data);
  }

  const importObject = {
    env: {
      // https://github.com/suborbital/reactr/tree/a6577851496b932f1f59bc58da46c74e1c0a175a/rwasm/api
      log_msg(pointer, size, level, identifier) {
        const str = readString(pointer, size);
        console.log(level, str);
      },
      cache_set(keyPointer, keySize, valPointer, valSize, ttl, identifier) {
        const key = readString(keyPointer, keySize);
        const val = readData(valPointer, valSize);
        console.debug("setting cache key", key, val);
        cache[key] = val;
        return 0;
      },
      cache_get(keyPointer, keySize, identifier) {
        const key = readString(keyPointer, keySize);
        console.debug("getting cache key", key);
        if (cache[key]) {
          ffiResult = cache[key];
          return ffiResult.length;
        }

        return 0;
      },
      get_ffi_result(pointer, identifier) {
        if (ffiResult == null) {
          return -1;
        }

        const result = ffiResult;
        ffiResult = null;
        // FIXME: no check on size?
        writeData(pointer, result);
      },
      request_get_field(fieldType, keyPointer, keySize, identifier) {
        const key = readString(keyPointer, keySize);
        console.debug("request_get_field", key);

        switch(fieldType) {
          case RequestFieldTypeMeta:
            switch (key) {
              case "method":
                ffiResult = method;
                return ffiResult.length
                break;
              default:
                throw new Error("unsupported RequestFieldTypeMeta: " + key)
            }
            break;
          case RequestFieldTypeParams:
            if (params[key]) {
              ffiResult = params[key];
              return ffiResult.length
            }
            break;
          default:
            throw new Error("unsupported request_get_field type: " + fieldType)
        }
      },
      return_result(pointer, size, ident) {
        const str = readString(pointer, size);
        response = new Response(str);
      },
      return_error(code, pointer, size, ident) {
        const str = readString(pointer, size);
        response = new Response(`Error ${code}: ${str}`, { status: 500 });
      }
    },
    wasi_snapshot_preview1: {
      fd_write() {
        console.log("fd_write", arguments);
      },
      proc_exit() {
        console.log("proc_exit", arguments);
      },
      environ_sizes_get() {
        console.log("environ_sizes_get", arguments);
        return 0;
      },
      environ_get() {
        console.log("environ_get", arguments);
      },
    }
  };

  const instance = new WebAssembly.Instance(WASM_MODULE, importObject);
  exports = instance.exports;

  const inPtr = allocate(jobBytes.length);
  writeData(inPtr, jobBytes)
  jobParams = params;

  if (typeof exports.init === "function") {
    exports.init();
  }
  if (typeof exports._start === "function") {
    exports._start();
  }
  exports.run_e(inPtr, jobBytes.length, ident);

  // FIXME: deallocate jobbytes

  return response;
}
