package services

var clientServingMode string

func InitServingMode(mode string) {
	clientServingMode = mode
}

func IsAsyncMode() bool {
	return clientServingMode == "async"
}
