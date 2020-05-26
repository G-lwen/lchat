package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"lchat/service/conf"
	"lchat/service/entity"
	"net/http"
)

var githubState = uuid.NewV4().String()

// token
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
}

// github 返回的用户信息
type GithubUserInfo struct {
	Login string `json:"login"`
	ID uint64 `json:"id"`
    NodeId string `json:"node_id"`
	AvatarUrl string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	Url string `json:"url"`
	HtmlUrl string `json:"html_url"`
	FollowersUrl string `json:"followers_url"`
	FollowingUrl string `json:"following_url"`
	GistsUrl string `json:"gists_url"`
	StarredUrl string `json:"starred_url"`
	SubscriptionsUrl string `json:"subscriptions_url"`
	OrganizationsUrl string `json:"organizations_url"`
	ReposUrl string `json:"repos_url"`
	EventsUrl string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type string  `json:"type"`
	SiteAdmin bool  `json:"site_admin"`
	Name interface{}  `json:"name"`
	Company interface{}  `json:"company"`
	Blog string  `json:"blog"`
	Location interface{}  `json:"location"`
	Email interface{}  `json:"email"`
	Hireable interface{}  `json:"hireable"`
	Bio interface{}  `json:"bio"`
	PublicRepos uint64  `json:"public_repos"`
	PublicGists uint64  `json:"public_gists"`
	Followers uint64  `json:"followers"`
	Following uint64  `json:"following"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// 向第三方发送请求授权
func oauth(c *gin.Context) {
	oauthType := c.Param("oauthType")
	url := "/login"
	if oauthType == "github" {
		url = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_url=%s&scope=user:email&state=%s",
			conf.Get().Oauth.GithubClientId, conf.Get().Oauth.GithubRedirectUrl, githubState)
	}
	c.Redirect(http.StatusFound, url)
}

// 第三方请求授权回调
func oauthCallback(c *gin.Context) {
	url := "/login"
	oauthType := c.Param("oauthType")
	code := c.Query("code")
	state := c.Query("state")

	if oauthType == "github" {
		// 校验 state 是否正确
		if state == "" || state != githubState {
			c.Redirect(http.StatusFound, url)
			return
		}

		// 请求获取 token
		token, err := getGithubToken(code)
		if err != nil {
			c.Redirect(http.StatusFound, url)
			return
		}
		// 通过 token 获取用户信息
		userInfo, err := getUserInfoByGithubToken(token)
		if err != nil {
			c.Redirect(http.StatusFound, url)
			return
		}
        user := &entity.User{
        	GithubOpenId: userInfo.ID,
		}

        if err = user.LoadByGithubOpenId(); err != nil || user.ID == 0 {
			// 用户第一次使用Github登录，先进行登记
			user.GithubUrl = userInfo.HtmlUrl
			user.AvatarUrl = userInfo.AvatarUrl
			user.NickName = userInfo.Login
			if err := user.Save(); err != nil {
				c.Redirect(http.StatusFound, url)
				return
			}
		}
		session := sessions.Default(c)
		session.Clear()
		session.Set(sessionKey, user)
		session.Save()
		url = "/"
	}
	c.Redirect(http.StatusFound, url)
}

// 获取 github 的登录 token
func getGithubToken(code string) (*Token, error) {
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		conf.Get().Oauth.GithubClientId, conf.Get().Oauth.GithubClientSecret, code)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	token := &Token{}
	if err = json.NewDecoder(res.Body).Decode(token); err != nil {
		return nil, err
	}
	return token, nil
}

// 通过 github token 获取需要的用户信息
func getUserInfoByGithubToken(token *Token) (*GithubUserInfo, error) {
	url := fmt.Sprintf("https://api.github.com/user?access_token=%s", token.AccessToken)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	userInfo := &GithubUserInfo{}
	if err = json.NewDecoder(res.Body).Decode(userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}
