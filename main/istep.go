package main

type istep interface {
	getFilePath() string
	getFileCode() string
}

