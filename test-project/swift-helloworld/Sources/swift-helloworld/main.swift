import Suborbital

class SwiftHelloworldRunnable: Suborbital.Runnable {
    func run(input: String) -> String {
        return "hello " + input
    }
}

@_cdecl("init")
func `init`() {
    Suborbital.Set(runnable: SwiftHelloworldRunnable())
}