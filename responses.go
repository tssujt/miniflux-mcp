package main

import (
	"encoding/json"

	"miniflux.app/v2/client"
)

func summarizeFeed(feed *client.Feed) *feedSummary {
	if feed == nil {
		return nil
	}
	summary := &feedSummary{
		ID: feed.ID, Title: feed.Title, FeedURL: feed.FeedURL, Disabled: feed.Disabled,
	}
	if feed.Category != nil {
		summary.Category = &feedCategorySummary{ID: feed.Category.ID, Title: feed.Category.Title}
	}
	return summary
}

// Embed the original entry to preserve article fields, overriding only its feed.
type entryWithCompactFeed struct {
	*client.Entry
	Feed *feedSummary `json:"feed,omitempty"`
}

func marshalEntryList(result *client.EntryResultSet) ([]byte, error) {
	if result == nil {
		return json.MarshalIndent(result, "", "  ")
	}
	var entries []*entryWithCompactFeed
	if result.Entries != nil {
		entries = make([]*entryWithCompactFeed, len(result.Entries))
		for i, entry := range result.Entries {
			if entry != nil {
				entries[i] = &entryWithCompactFeed{Entry: entry, Feed: summarizeFeed(entry.Feed)}
			}
		}
	}
	return json.MarshalIndent(struct {
		*client.EntryResultSet
		Entries []*entryWithCompactFeed `json:"entries"`
	}{EntryResultSet: result, Entries: entries}, "", "  ")
}
