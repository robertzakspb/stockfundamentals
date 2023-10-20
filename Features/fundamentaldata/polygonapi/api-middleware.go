package polygonapi

import (
	"time"
)

const APIKey = "kATfwcnbHHRPIxQgOKRF0iFYT0hWx62i"

// Logs the timestamps of requests to Polygon API
var requestTimeStamps []time.Time

// Polygon API has a 5-requests-per-minute limit for users of their free tier
// This method ensures that we stay within the limit
func canSendRequestToPolygon() bool {
	if len(requestTimeStamps) < 5 {
		return true
	}

	oneMinuteAgo := time.Now().Add(-61 * time.Second)
	fifthFromLastRequestTimestamp := requestTimeStamps[len(requestTimeStamps)-5]

	return fifthFromLastRequestTimestamp.Before(oneMinuteAgo)
}

func DelayRequestIfAPILimitReached() {
	if !canSendRequestToPolygon() {
		time.Sleep(61 * time.Second)
	}
	requestTimeStamps = append(requestTimeStamps, time.Now())
}
