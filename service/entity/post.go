package entity

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
	"lchat/service/utils"
)

// 文章
type Post struct {
	Entity

	Title string                                          // 文章标题
	Summary string                                        // 文章摘要
	Body string `gorm:"type:text"`                        // 文章 默认 markdown 内容
	HtmlBody string `gorm: "type:text"`                   // 文章 html 内容
	UserId uint                                           // 文章作者ID
	View uint                                             // 浏览量
	Published bool                                       // 是否发布
	Tags []*Tag `gorm:"-"`                                // 文章标签
}

// 标签
type Tag struct {
	Entity

	Name string `gorm:"unique_index"`                      // 标签名称
	Total int `gorm:"-"`                                   // 标签文章量
}

// 文章标签
type PostTag struct {
	PostId uint `gorm:"primary_key;auto_increment:false"`   // 文章 ID
	TagId uint `gorm:"primary_key;auto_increment:false"`    // 标签 ID
}


//////////////////////////       Post       ///////////////////////////////

// 保存文章
func (post *Post) Save() error {
	return store[PostStore].Create(post).Error
}

// 更新文章内容
func (post *Post) Update() error {
	return store[PostStore].Model(post).Updates(map[string]interface{}{
		"title": post.Title,
		"body": post.Body,
		"summary": post.Summary,
		"html_body": post.HtmlBody,
		"published": post.Published,
	}).Error
}

// 更新文章浏览量
func (post *Post) UpdateView() error {
	return store[PostStore].Model(post).Update("view", post.View).Error
}

// 根据ID删除文章
func (post *Post) Delete() error {
	return store[PostStore].Delete(post).Error
}

// 根据 ID 加载文章内容
func (post *Post) Load() error {
	return store[PostStore].First(post).Error
}

// 提取文章摘要
func (post *Post) ExtractSummary() {
	if post.Body == "" {
		return
	}
	unsafe := blackfriday.Run([]byte(post.Body))
	html := bluemonday.UGCPolicy().Sanitize(string(unsafe))
	post.Summary = utils.TrimHtml(html)
	if len(post.Summary) > 150 {
		post.Summary = post.Summary[:150] + "..."
	}
}


////////////////////////       PostTag       //////////////////////////////

// 添加文章标签
func (post *Post) AddTag(tagId uint) error {
	postTag := &PostTag{
		TagId: tagId,
		PostId: post.ID,
	}
	return store[PostStore].FirstOrCreate(postTag, "tag_id = ? and post_id = ?", postTag.TagId, postTag.PostId).Error
}

// 删除文章标签
func (post *Post) RemoveTag(tagId uint) error {
	return store[PostStore].Delete(PostTag{}, "tag_id = ? and post_id = ?", tagId, post.ID).Error
}

// 加载文章标签
func (post *Post) LoadTags() error {
	var tags []*Tag
	sql := `select t.* from tags t 
			inner join post_tags pt on t.id = pt.tag_id 
			where t.deleted_at is null and pt.post_id = ?`
	rows, err := store[PostStore].Raw(sql, post.ID).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		store[PostStore].ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	post.Tags = tags
	return nil
}

//////////////////////////       Tag       ////////////////////////////////

// 加载标签，如果标签不存在，则进行保存处理
func (tag *Tag) LoadByName() error {
	return store[PostStore].FirstOrCreate(tag, "name = ?", tag.Name).Error
}

////////////////////////       Service       //////////////////////////////

// 分页查询发布的文章
func ListPosts(pageIndex, pageSize int) ([]*Post, int, error) {
	return listPosts(pageIndex, pageSize, 0, true)
}

// 根据 文章作者ID 分页查询发布的文章
func ListPostsByUserId(pageIndex, pageSize int, userId uint, published bool) ([]*Post, int, error) {
	return listPosts(pageIndex, pageSize, userId, published)
}

// 多条件查询文章
func listPosts(pageIndex, pageSize int, userId uint, published bool) ([]*Post, int, error) {
	var (
		posts []*Post
		count int
	)
	db := store[PostStore].Where("published = ?", published)
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}
	err := db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&posts).Count(&count).Error
	return posts, count, err
}

// 根据 标签 分页查询发布的文章
func ListPostsByTag(pageIndex, pageSize int, tagId uint) ([]*Post, error) {
	var posts []*Post
	sql := `select p.* from posts p 
			inner join post_tags pt on p.id = pt.post_id 
			where p.deleted_at is null and p.published = ? and pt.tag_id = ? limit ?, ?`
	rows, err := store[PostStore].Raw(sql, true, tagId, (pageIndex - 1) * pageSize, pageSize).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		store[PostStore].ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return posts, err
}

// 根据 标签 获取发布的文章数
func CountPostsByTag(tagId uint) int {
	var count int
	sql := `select count(*) from posts p 
			inner join post_tags pt on p.id = pt.post_id 
			where p.deleted_at is null and p.published = ? and pt.tag_id = ?`
	store[PostStore].Raw(sql, true, tagId).Row().Scan(&count)
	return count
}

// 查询标签
func ListAllTags() ([]*Tag, error) {
	var tags []*Tag
	sql := `select t.*, count(pt.post_id) total from tags t 
		left join (
			select post_tags.* from post_tags 
			inner join posts on post_tags.post_id = posts.id 
			where posts.deleted_at is null and posts.published = ?
		) pt on t.id = pt.tag_id
		where t.deleted_at is null group by t.id`
	rows, err := store[PostStore].Raw(sql, true).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		store[PostStore].ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}

// 查询标签数
func CountTags() int {
	var count int
	store[PostStore].Model(&Tag{}).Count(&count)
	return count
}
