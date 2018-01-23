package server

// Init -- setup the server
func Init() {
	r := NewRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}
