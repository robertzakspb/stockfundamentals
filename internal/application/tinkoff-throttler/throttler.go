// Throttler for various Tinkoff APIs
package tthrottler

import (
	"time"
)

// Multiplication by 0.95 is just in case
const instrumentServiceLimit = time.Minute / (200 * 0.95)
// const accountServiceLimit = time.Minute / (100 * 0.95)
// const operationServiceLimit = time.Minute / (200 * 0.95)
const marketDataServiceLimit = time.Second //Not sure what is the right limit but this one works

var InstrumentServiceThrottle = time.Tick(instrumentServiceLimit)
var MarketDataServiceThrottle = time.Tick(marketDataServiceLimit)