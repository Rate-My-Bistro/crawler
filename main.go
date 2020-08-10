package main

import "github.com/ansgarS/rate-my-bistro-crawler/restapi"

// the application cycle
func main() {
	//TODO
	// - persister für jeden datentyp
	// - REST API POST das einen neuen job erstellt: datum im body nicht in der url
	restapi.Serve()
}
