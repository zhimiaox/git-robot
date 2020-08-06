package main

var Config = &cfg{}

type cfg struct {
	Server struct {
		APIListen    string
		ReadTimeOut  int
		WriteTimeOut int
	}
	MapStorage struct {
		FilePath string
	}
	Work struct {
		GitPath string
		MaxErrCount int
	}
}
