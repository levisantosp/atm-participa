package tests

func GetCookie(sessionID string) string {
	return "Cookie: session=" + sessionID
}
