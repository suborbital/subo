package greetings

// This comment is associated with the ExampleHello_doNotDoThis plugin example.
// This is showing you how not to use Hello.
func ExampleHello_doNotDoThis() {
	Hello("Boo!")
	// Output: Ahhh!
}

// This comment is associated with the ExampleHello_doThis plugin example.
// This is showing you how to use Hello.
func ExampleHello_doThis() {
	Hello("world")
	// Output: Hello, world!
}

// This comment is associated with the package example.
// This is showing you how to use Hello.
func Example() {
	Hello("Foo Bar")
	// Output: Hello, Foo Bar!
}
