import Suborbital

class SuboRunnable: Suborbital.Runnable {
    func run(input: String) -> String {
        return "subo says hello " + input
    }
}

@_cdecl("init")
func `init`() {
    Suborbital.Set(runnable: SuboRunnable())
}