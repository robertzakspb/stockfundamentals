// Throttler for various Tinkoff APIs
package tthrottler

import (
	"time"
)

// Multiplication by 0.95 is just in case
const instrumentServiceLimit = time.Minute / (200 * 0.95)
const marketDataServiceLimit = time.Minute / (600 * 0.95) //Not sure what is the right limit but this one works
// const accountServiceLimit = time.Minute / (100 * 0.95)
// const operationServiceLimit = time.Minute / (200 * 0.95)

var InstrumentServiceThrottle = time.Tick(instrumentServiceLimit)
var MarketDataServiceThrottle = time.Tick(marketDataServiceLimit)