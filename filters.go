package main

import "miniflux.app/v2/client"

func newEntryFilter() *client.Filter {
	return &client.Filter{Limit: -1, Offset: -1}
}

func parseEntryFilter(argsMap map[string]interface{}) *client.Filter {
	filter := newEntryFilter()

	if statusValues, ok := argsMap["statuses"].([]interface{}); ok {
		for _, statusValue := range statusValues {
			if status, ok := statusValue.(string); ok {
				filter.Statuses = append(filter.Statuses, status)
			}
		}
	}

	if statusStr, ok := argsMap["status"].(string); ok && len(filter.Statuses) == 0 {
		filter.Status = statusStr
	}

	int64Fields := map[string]*int64{
		"category_id":      &filter.CategoryID,
		"after":            &filter.After,
		"before":           &filter.Before,
		"published_after":  &filter.PublishedAfter,
		"published_before": &filter.PublishedBefore,
		"changed_after":    &filter.ChangedAfter,
		"changed_before":   &filter.ChangedBefore,
		"after_entry_id":   &filter.AfterEntryID,
		"before_entry_id":  &filter.BeforeEntryID,
	}
	for name, target := range int64Fields {
		if value, ok := argsMap[name].(float64); ok {
			*target = int64(value)
		}
	}

	intFields := map[string]*int{
		"limit":  &filter.Limit,
		"offset": &filter.Offset,
	}
	for name, target := range intFields {
		if value, ok := argsMap[name].(float64); ok {
			*target = int(value)
		}
	}

	stringFields := map[string]*string{
		"search":    &filter.Search,
		"order":     &filter.Order,
		"direction": &filter.Direction,
	}
	for name, target := range stringFields {
		if value, ok := argsMap[name].(string); ok {
			*target = value
		}
	}

	if starred, ok := argsMap["starred"].(bool); ok {
		if starred {
			filter.Starred = client.FilterOnlyStarred
		} else {
			filter.Starred = client.FilterNotStarred
		}
	}

	if globallyVisible, ok := argsMap["globally_visible"].(bool); ok {
		filter.GloballyVisible = globallyVisible
	}

	return filter
}
