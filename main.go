package main

func main() {
	getConfig()
	initLoggers()
	go delayServer()
	go downloadServer()
	go indexingClient()

	RunApplication()
}
