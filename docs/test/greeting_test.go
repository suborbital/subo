package greetings

// This comment is associated with the ExampleHello_doNotDoThis function example.
// This is showing you how not to use Hello
func ExampleHello_doNotDoThis() {
	Hello("Boo!")
	// Output: Ahhh!
}

// This comment is associated with the ExampleHello_doThis function example.
// This is showing you how to use Hello
func ExampleHello_doThis() {
	Hello("world")
	// Output: Hello, world!
}

// This comment is associated with the package example.
// This is showing you how to use Hello
func Example() {
	Hello("Foo Bar")
	// Output: Hello, Foo Bar!
}

// This comment is associated with the package example Example_doNotDoThis.
// This is showing you how to use Hello
func Example_doNotDoThis() {
	Hello("boo!")
	// Output: Hello, Bar Foo!
}
