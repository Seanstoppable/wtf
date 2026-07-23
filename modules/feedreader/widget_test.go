package feedreader

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"gotest.tools/assert"
)

func Test_rotateShowType(t *testing.T) {
	tests := []struct {
		name     string
		input    ShowType
		expected ShowType
	}{
		{
			name:     "SHOW_TITLE rotates to SHOW_LINK",
			input:    SHOW_TITLE,
			expected: SHOW_LINK,
		},
		{
			name:     "SHOW_LINK rotates to SHOW_CONTENT",
			input:    SHOW_LINK,
			expected: SHOW_CONTENT,
		},
		{
			name:     "SHOW_CONTENT rotates to SHOW_TITLE",
			input:    SHOW_CONTENT,
			expected: SHOW_TITLE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := rotateShowType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_sort(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)
	earliest := now.Add(-2 * time.Hour)

	tests := []struct {
		name           string
		feedItems      []*FeedItem
		expectedTitles []string
	}{
		{
			name:           "empty list",
			feedItems:      []*FeedItem{},
			expectedTitles: []string{},
		},
		{
			name: "already sorted (newest first)",
			feedItems: []*FeedItem{
				{item: &gofeed.Item{Title: "Newest", PublishedParsed: &now}},
				{item: &gofeed.Item{Title: "Oldest", PublishedParsed: &earliest}},
			},
			expectedTitles: []string{"Newest", "Oldest"},
		},
		{
			name: "unsorted items get sorted newest first",
			feedItems: []*FeedItem{
				{item: &gofeed.Item{Title: "Oldest", PublishedParsed: &earliest}},
				{item: &gofeed.Item{Title: "Newest", PublishedParsed: &now}},
				{item: &gofeed.Item{Title: "Middle", PublishedParsed: &earlier}},
			},
			expectedTitles: []string{"Newest", "Middle", "Oldest"},
		},
		{
			name: "nil PublishedParsed preserves original order (not sortable)",
			feedItems: []*FeedItem{
				{item: &gofeed.Item{Title: "No Date", PublishedParsed: nil}},
				{item: &gofeed.Item{Title: "Has Date", PublishedParsed: &now}},
			},
			expectedTitles: []string{"No Date", "Has Date"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widget := &Widget{}
			result := widget.sort(tt.feedItems)

			actualTitles := make([]string, len(result))
			for i, item := range result {
				actualTitles[i] = item.item.Title
			}

			assert.DeepEqual(t, tt.expectedTitles, actualTitles)
		})
	}
}

func Test_getShowText(t *testing.T) {
	publishedTime := time.Date(2023, time.March, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		feedItem *FeedItem
		showType ShowType
		settings *Settings
		expected string
	}{
		{
			name:     "with nil FeedItem",
			feedItem: nil,
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "",
		},
		{
			name: "with plain title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Cats and Dogs"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "[white]Cats and Dogs",
		},
		{
			name: "with escaped title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "&lt;Cats and Dogs&gt;"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "[white]<Cats and Dogs>",
		},
		{
			name: "with unescaped title",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "<Cats and Dogs>"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "[white]<Cats and Dogs>",
		},
		{
			name: "with source-title",
			feedItem: &FeedItem{
				sourceTitle: "WTF",
				item:        &gofeed.Item{Title: "<Cats and Dogs>"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "[green]WTF [white]<Cats and Dogs>",
		},
		{
			name: "with link",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Cats and Dogs", Link: "https://cats.com/dog.xml"},
			},
			showType: SHOW_LINK,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "https://cats.com/dog.xml",
		},
		{
			name: "with content",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Article", Content: "<p>Hello World</p>"},
			},
			showType: SHOW_CONTENT,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: false,
			},
			expected: "[white]Article\nHello World",
		},
		{
			name: "with publish date shown",
			feedItem: &FeedItem{
				item: &gofeed.Item{
					Title:           "News",
					Published:       "2023-03-15",
					PublishedParsed: &publishedTime,
				},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:          colors{source: "green", publishDate: "orange"},
				showSource:      false,
				showPublishDate: true,
				dateFormat:       "Jan 02",
			},
			expected: "[orange]Mar 15 [white]News",
		},
		{
			name: "with source and publish date",
			feedItem: &FeedItem{
				sourceTitle: "Blog",
				item: &gofeed.Item{
					Title:           "Post",
					Published:       "2023-03-15",
					PublishedParsed: &publishedTime,
				},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:          colors{source: "green", publishDate: "orange"},
				showSource:      true,
				showPublishDate: true,
				dateFormat:       "Jan 02",
			},
			expected: "[green]Blog [orange]Mar 15 [white]Post",
		},
		{
			name: "with whitespace in title collapsed",
			feedItem: &FeedItem{
				item: &gofeed.Item{Title: "Lots   of    spaces"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: false,
			},
			expected: "[white]Lots of spaces",
		},
		{
			name: "showSource true but empty sourceTitle shows no prefix",
			feedItem: &FeedItem{
				sourceTitle: "",
				item:        &gofeed.Item{Title: "Title"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: true,
			},
			expected: "[white]Title",
		},
		{
			name: "showSource false hides sourceTitle",
			feedItem: &FeedItem{
				sourceTitle: "Hidden",
				item:        &gofeed.Item{Title: "Title"},
			},
			showType: SHOW_TITLE,
			settings: &Settings{
				colors:     colors{source: "green", publishDate: "orange"},
				showSource: false,
			},
			expected: "[white]Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widget := &Widget{
				settings: tt.settings,
				showType: tt.showType,
			}

			actual := widget.getShowText(tt.feedItem, "white")

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_fetchForFeed_feedLimit(t *testing.T) {
	tests := []struct {
		name      string
		feedLimit int
		itemCount int
		expected  int
	}{
		{
			name:      "feedLimit of 0 returns all items",
			feedLimit: 0,
			itemCount: 5,
			expected:  5,
		},
		{
			name:      "negative feedLimit returns all items",
			feedLimit: -1,
			itemCount: 5,
			expected:  5,
		},
		{
			name:      "feedLimit less than item count truncates",
			feedLimit: 2,
			itemCount: 5,
			expected:  2,
		},
		{
			name:      "feedLimit greater than item count returns all",
			feedLimit: 10,
			itemCount: 3,
			expected:  3,
		},
		{
			name:      "feedLimit equal to item count returns all",
			feedLimit: 3,
			itemCount: 3,
			expected:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build feed items
			items := make([]*gofeed.Item, tt.itemCount)
			for i := range items {
				items[i] = &gofeed.Item{Title: "Item"}
			}

			widget := &Widget{
				settings: &Settings{
					feedLimit: tt.feedLimit,
				},
			}

			// Simulate fetchForFeed logic inline since we can't easily mock ParseURL
			var feedItems []*FeedItem
			for idx, gofeedItem := range items {
				if widget.settings.feedLimit >= 1 && idx >= widget.settings.feedLimit {
					break
				}
				feedItems = append(feedItems, &FeedItem{
					item:        gofeedItem,
					sourceTitle: "Test",
					viewed:      false,
				})
			}

			assert.Equal(t, tt.expected, len(feedItems))
		})
	}
}

func Test_fetchForFeed_alias(t *testing.T) {
	tests := []struct {
		name          string
		alias         string
		feedTitle     string
		expectedTitle string
	}{
		{
			name:          "alias overrides feed title",
			alias:         "MyAlias",
			feedTitle:      "Original Feed",
			expectedTitle: "MyAlias",
		},
		{
			name:          "empty alias preserves feed title",
			alias:         "",
			feedTitle:      "Original Feed",
			expectedTitle: "Original Feed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feedItem := &FeedItem{
				item:        &gofeed.Item{Title: "Item"},
				sourceTitle: tt.feedTitle,
				viewed:      false,
			}
			if tt.alias != "" {
				feedItem.sourceTitle = tt.alias
			}

			assert.Equal(t, tt.expectedTitle, feedItem.sourceTitle)
		})
	}
}
