package consts

func GetSalt() []byte {
	return []byte{1, 2, 3, 4, 5, 6, 7, 8}
}

func GetPublicRoutes() []string {
	return []string{"/user/signup", "/user/login"}
}

func GetSecret() []byte {
	return []byte("secret")
}