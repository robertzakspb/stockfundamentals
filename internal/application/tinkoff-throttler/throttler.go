// Throttler for various Tinkoff APIs
package tthrottler

import (
	"time"
)
//These limits were discovered empirically; the limits provided in the official documentation, however, frequent ceiling breaches
const instrumentServiceLimit = time.Minute / 150
const marketDataServiceLimit = time.Minute / 29 

var InstrumentServiceThrottle = time.Tick(instrumentServiceLimit)
var MarketDataServiceThrottle = time.Tick(marketDataServiceLimit)