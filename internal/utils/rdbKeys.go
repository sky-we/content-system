package utils

import "fmt"

func GenSessionKey(userId string) string {
	return fmt.Sprintf("session_id:%s", userId)
}

func GenAuthKey(sessionId string) string {
	return fmt.Sprintf("session_auth:%s", sessionId)
}
