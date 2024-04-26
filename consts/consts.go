package consts

func GetSalt() []byte {
	return []byte{1, 2, 3, 4, 5, 6, 7, 8}
}

func GetArgonOptions() ([]byte, uint32, uint32, uint8, uint32) {
	return GetSalt(), 3, 32 * 1024, 4, 32
}

func GetPublicRoutes() []string {
	return []string{"/user/signup", "/user/login"}
}

func GetSecret() []byte {
	return []byte("secret")
}

func GetPort() string {
	return "3333"
}