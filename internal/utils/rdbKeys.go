package utils

import "fmt"

const SessionKey = "session_id"

func GenSessionKey(userId string) string {
	return fmt.Sprintf("session_id:%s", userId)
}

func GenAuthKey(sessionId string) string {
	return fmt.Sprintf("session_auth:%s", sessionId)
}
