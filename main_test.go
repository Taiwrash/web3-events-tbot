package main

import (
	"testing"

	"github.com/Taiwrash/web3event-spot/scrape"
	"github.com/stretchr/testify/assert"
)

func TestScrape(t *testing.T) {
	events, err := scrape.Scrape()
	assert.NoError(t, err)
	assert.Equal(t, 19, len(events))

	for i, event := range events {
		assert.NotEmpty(t, event.Title, "Event %d has an empty title", i+1)
		assert.NotEmpty(t, event.Location, "Event %d has an empty location", i+1)
		assert.NotEmpty(t, event.Date, "Event %d has an empty date", i+1)
		assert.NotEmpty(t, event.URL, "Event %d has an empty URL", i+1)
	}
}
