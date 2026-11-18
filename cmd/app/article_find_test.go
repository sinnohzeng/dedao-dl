package app

import (
	"errors"
	"testing"

	"github.com/yann0917/dedao-dl/services"
)

// pages serves ids newest-first, 3 per page, honouring the max_id cursor unless stuck is set.
func pages(ids []int, stuck bool, calls *int) func(int) (*services.ArticleList, error) {
	return func(maxID int) (*services.ArticleList, error) {
		*calls++
		start := 0
		if maxID != 0 && !stuck {
			for i, id := range ids {
				if id == maxID {
					start = i + 1
				}
			}
		}
		end := start + 3
		if end > len(ids) {
			end = len(ids)
		}
		list := &services.ArticleList{}
		for _, id := range ids[start:end] {
			a := services.ArticleIntro{}
			a.ID = id
			list.List = append(list.List, a)
		}
		return list, nil
	}
}

func TestFindArticle(t *testing.T) {
	ids := []int{90, 80, 70, 60, 50, 40, 30}
	cases := []struct {
		name  string
		aid   int
		stuck bool
		found bool
		calls int
	}{
		{"newest lecture costs one page", 90, false, true, 1},
		{"older lecture keeps paging", 40, false, true, 2},
		{"missing lecture stops at the empty page", 99, false, false, 4},
		{"cursor ignored stops after a page of seen ids", 30, true, false, 2},
	}
	for _, c := range cases {
		calls := 0
		got, err := findArticle(pages(ids, c.stuck, &calls), c.aid)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if (got != nil) != c.found || calls != c.calls {
			t.Fatalf("%s: found=%v calls=%d", c.name, got != nil, calls)
		}
		if got != nil && (len(got.List) != 1 || got.List[0].ID != c.aid) {
			t.Fatalf("%s: list %+v", c.name, got.List)
		}
	}
}

func TestFindArticleStopsOnError(t *testing.T) {
	calls := 0
	_, err := findArticle(func(int) (*services.ArticleList, error) {
		calls++
		return nil, errors.New("HTTP 405")
	}, 1)
	if err == nil || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}
