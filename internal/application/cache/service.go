package cache

import "time"

func LoadAllCache() {
	go LoadForexCache()
	go LoadQuoteCache()
}

// Starts a perpetual loop that updates different caches with the required frequency
func StartCacheUpdateLoop() {
	//Initial load
	LoadAllCache()

	//Looping the updates
	go startQuoteCacheUpdateLoop()

	//Forex rates are updated daily and hence do not require constant cache updating beyond the initial load
}

func startQuoteCacheUpdateLoop() {
	var quoteCacheUpdateThrottler = time.Tick(time.Minute)
	for {
		<-quoteCacheUpdateThrottler
		LoadQuoteCache()
	}
}
