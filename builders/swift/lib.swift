

@_silgen_name("return_result")
func return_result(result_pointer: UnsafeRawPointer, result_size: Int32, ident: Int32)
@_silgen_name("print")
func print_msg(pointer: UnsafeRawPointer, size: Int32, ident: Int32)

@_cdecl("run_e")
func run_e(pointer: UnsafeRawPointer, size: Int32, ident: Int32) {
    let typed: UnsafePointer<UInt8> = pointer.bindMemory(to: UInt8.self, capacity: Int(size))
    let inString = String(cString: typed)
    
    let printMsg = "testing!"
    let printPtr = UnsafeRawPointer(printMsg)
    let printCount = Int32(printMsg.utf8.count)
    print_msg(pointer: printPtr, size: printCount, ident: ident)
    
    let retString = "hello " + inString

    let count = Int32(retString.utf8.count)

    let retPointer = UnsafeRawPointer(retString)

    return_result(result_pointer: retPointer, result_size: count, ident: ident)
}

@_cdecl("allocate")
func allocate(size: Int) -> UnsafeMutableRawPointer {
  return UnsafeMutableRawPointer.allocate(byteCount: size, alignment: MemoryLayout<UInt8>.alignment)
}

@_cdecl("deallocate")
func deallocate(pointer: UnsafeMutableRawPointer, size: Int) {
  pointer.deallocate()
}
