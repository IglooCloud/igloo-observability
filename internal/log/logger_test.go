package log

func ExampleLogger() {
	logger := Default()

	logger.Print("Test")
}
