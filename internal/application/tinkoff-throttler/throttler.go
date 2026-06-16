// Throttler for various Tinkoff APIs
package tthrottler

import (
	"time"
)

// Multiplication by 0.95 is just in case
const instrumentServiceLimit = time.Minute / 195
const marketDataServiceLimit = time.Minute / 29 //Not sure what is the right limit but this one works

var InstrumentServiceThrottle = time.Tick(instrumentServiceLimit)
var MarketDataServiceThrottle = time.Tick(marketDataServiceLimit)