package app

import (
	"errors"
	"fmt"

	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// ArticleList 已购课程文章列表
func ArticleList(id int, chapterID string) (list *services.ArticleList, err error) {
	info, err := CourseDetail(CateCourse, id)
	if err != nil {
		return
	}
	enid := info.Enid
	return ArticleListAllByEnId(enid, chapterID)
}

func ArticleListAllByEnId(enid, chapterID string) (list *services.ArticleList, err error) {
	maxID := 0
	var lists []services.ArticleIntro
	for {
		pageList, err1 := getService().ArticleList(enid, chapterID, maxID)
		if err1 != nil {
			err = err1
			return
		}
		if pageList == nil || len(pageList.List) == 0 {
			break
		}
		lists = append(lists, pageList.List...)
		lastID := pageList.List[len(pageList.List)-1].ID
		if lastID == maxID {
			break
		}
		maxID = lastID
		list = pageList
	}
	if list == nil {
		list = &services.ArticleList{}
	}
	list.List = lists
	return
}

// courseArticles lists the whole course, or only the lecture aid when aid > 0.
func courseArticles(enid string, aid int) (*services.ArticleList, error) {
	if aid <= 0 {
		return ArticleListAllByEnId(enid, "")
	}
	// New lectures sit at the end of the list, so walk it newest-first and stop at the match instead of
	// listing the whole course for one lecture. A miss falls back to the full forward list.
	hit, err := findArticle(func(maxID int) (*services.ArticleList, error) {
		return getService().ArticleListPage(enid, "", maxID, true)
	}, aid)
	if err != nil || hit != nil {
		return hit, err
	}
	return ArticleListAllByEnId(enid, "")
}

// findArticle pages with the max_id cursor until aid shows up, a page is empty, or a page brings no unseen id.
func findArticle(page func(maxID int) (*services.ArticleList, error), aid int) (*services.ArticleList, error) {
	seen := map[int]bool{}
	maxID := 0
	for {
		list, err := page(maxID)
		if err != nil {
			return nil, err
		}
		if list == nil || len(list.List) == 0 {
			return nil, nil
		}
		fresh := false
		for _, a := range list.List {
			if a.ID == aid {
				list.List = []services.ArticleIntro{a}
				return list, nil
			}
			if !seen[a.ID] {
				seen[a.ID] = true
				fresh = true
			}
		}
		if !fresh {
			return nil, nil
		}
		maxID = list.List[len(list.List)-1].ID
	}
}

func ArticleDetailByEnId(articleEnid string) (detail *services.ArticleDetail, err error) {
	info, err := getService().ArticleInfo(articleEnid, 1)
	if err != nil {
		return
	}
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, appid)
	return
}

// ArticleInfo article info
// get article token, audio token, media security token etc
func ArticleInfo(id, aid int) (info *services.ArticleInfo, aEnid string, err error) {
	list, err := ArticleList(id, "")
	if err != nil {
		return
	}

	var aids []int

	// get article enid
	for _, p := range list.List {
		aids = append(aids, p.ID)
		if p.ClassID == id && p.ID == aid {
			aEnid = p.Enid
		}
	}
	if !utils.Contains(aids, aid) {
		err = errors.New("找不到该文章 ID，请检查输入是否正确")
		return
	}

	info, err = getService().ArticleInfo(aEnid, 1)
	if err != nil {
		return
	}
	return
}

// ArticleDetail article detail
func ArticleDetail(id, aid int) (detail *services.ArticleDetail, enId string, err error) {
	info, aEnid, err := ArticleInfo(id, aid)
	if err != nil {
		return
	}
	enId = aEnid
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}

// ArticleCommentList article comment list
func ArticleCommentList(enId, sort string, page, limit int) (list *services.ArticleCommentList, err error) {
	list, err = getService().ArticleCommentList(enId, sort, page, limit)
	if err != nil {
		return
	}

	return

}

// OdobArticleInfo article info
// get article token, audio token, media security token etc.
func OdobArticleInfo(aEnid string) (info *services.ArticleInfo, err error) {
	info, err = getService().ArticleInfo(aEnid, 2)
	if err != nil {
		return
	}
	return
}

// OdobArticleDetail odob article detail
func OdobArticleDetail(aEnid string) (detail *services.ArticleDetail, err error) {
	info, err := OdobArticleInfo(aEnid)
	if err != nil {
		return
	}

	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}
