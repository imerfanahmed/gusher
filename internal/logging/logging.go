package logging

import "log"

// LogEvent logs event details
func LogEvent(appID, channel, event, socketID string) {
	log.Printf("Event: %s, Channel: %s, AppID: %s, SocketID: %s", event, channel, appID, socketID)
}