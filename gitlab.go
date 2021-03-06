package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var GitEmojiMap = map[string]string{
	":bulb:":                      "๐ก",
	":heavy_minus_sign:":          "โ",
	":bug:":                       "๐",
	":art:":                       "๐จ",
	":hammer:":                    "๐จ",
	":sparkles:":                  "โจ",
	":building_construction:":     "๐๏ธ",
	":wrench:":                    "๐ง",
	":triangular_flag_on_post:":   "๐ฉ",
	":arrow_down:":                "โฌ๏ธ",
	":label:":                     "๐ท๏ธ",
	":dizzy:":                     "๐ซ",
	":white_check_mark:":          "โ",
	":mag:":                       "๐๏ธ",
	":bento:":                     "๐ฑ",
	":chart_with_upwards_trend:":  "๐",
	":beers:":                     "๐ป",
	":boom:":                      "๐ฅ",
	":bookmark:":                  "๐",
	":monocle_face:":              "๐ง",
	":recycle:":                   "โป๏ธ",
	":card_file_box:":             "๐๏ธ",
	":globe_with_meridians:":      "๐",
	":adhesive_bandage:":          "๐ฉน",
	":pushpin:":                   "๐",
	":iphone:":                    "๐ฑ",
	":test_tube:":                 "๐งช",
	":page_facing_up:":            "๐",
	":alien:":                     "๐ฝ๏ธ",
	":children_crossing:":         "๐ธ",
	":poop:":                      "๐ฉ",
	":heavy_plus_sign:":           "โ",
	":necktie:":                   "๐",
	":rotating_light:":            "๐จ",
	":memo:":                      "๐",
	":loud_sound:":                "๐",
	":construction:":              "๐ง",
	":fire:":                      "๐ฅ",
	":zap:":                       "โก๏ธ",
	":stethoscope:":               "๐ฉบ",
	":package:":                   "๐ฆ๏ธ",
	":camera_flash:":              "๐ธ",
	":lipstick:":                  "๐",
	":mute:":                      "๐",
	":rocket:":                    "๐",
	":lock:":                      "๐๏ธ",
	":ambulance:":                 "๐๏ธ",
	":pencil2:":                   "โ๏ธ",
	":arrow_up:":                  "โฌ๏ธ",
	":clown_face:":                "๐คก",
	":truck:":                     "๐",
	":goal_net:":                  "๐ฅ",
	":egg:":                       "๐ฅ",
	":speech_balloon:":            "๐ฌ",
	":construction_worker:":       "๐ท",
	":passport_control:":          "๐",
	":rewind:":                    "โช๏ธ",
	":wheelchair:":                "โฟ๏ธ",
	":alembic:":                   "โ๏ธ",
	":seedling:":                  "๐ฑ",
	":green_heart:":               "๐",
	":tada:":                      "๐",
	":busts_in_silhouette:":       "๐ฅ",
	":twisted_rightwards_arrows:": "๐",
	":wastebasket:":               "๐๏ธ",
	":coffin:":                    "โฐ๏ธ",
	":see_no_evil:":               "๐",
}

func trans2Emoji(content string) string {
	for k, v := range GitEmojiMap {
		content = strings.ReplaceAll(content, k, v)
	}
	return content
}

func NewClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

type WxResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Push events
type PushBody struct {
	ObjectKind string     `json:"object_kind"`
	Ref        string     `json:"ref"`
	Commits    []Commit   `json:commits`
	Repository Repository `json:"repository"`
	After      string     `json:"after"`
	UserName   string     `json:"user_name"`
}

// TagPushBody Tag events
type TagPushBody struct {
	UserName   string     `json:"user_name"`
	Ref        string     `json:"ref"`
	Repository Repository `json:"repository"`
}

// IssuePushBody Issues events
type IssuePushBody struct {
	User             IssueUser   `json:"user"`
	Repository       Repository  `json:"repository"`
	ObjectAttributes IssueObject `json:"object_attributes"`
}

// CommentPushBody comment
type CommentPushBody struct {
	User             IssueUser     `json:"user"`
	Repository       Repository    `json:"repository"`
	ObjectAttributes CommentObject `json:"object_attributes"`
}

type CommentObject struct {
	Id        int64  `json:"id"`
	Note      string `json:"note"`
	UpdatedAt string `json:"updated_at"`
	Url       string `json:"url"`
}

// MRPushBody
type MRPushBody struct {
	User             IssueUser  `json:"user"`
	Repository       Repository `json:"repository"`
	ObjectAttributes MRObjects  `json:"object_attributes"`
}

// PipelineBody
type PipelineBody struct {
	ObjectAttributes PipelineObject `json:"object_attributes"`
	User             IssueUser      `json:"user"`
	Project          Project        `json:"project"`
}

type PipelineObject struct {
	Id         int64  `json:"id"`
	Ref        string `json:"ref"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	FinishedAt string `json:"finished_at"`
	Duration   int64  `json:"duration"`
	Tag        bool   `json:"tag"`
}

type MRObjects struct {
	Id           int64  `json:"id"`
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
	UpdatedAt    string `json:"updated_at"`
	Url          string `json:"url"`
	Action       string `json:"action"`
}

type IssueUser struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
}

type IssueObject struct {
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Url    string `jso:"url"`
	Action string `json:"action"`
}

type Commit struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	TimeStamp string `json:"timestamp"`
	Url       string `json:"url"`
	Author    Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository struct {
	Name      string `json:"name"`
	HomePage  string `json:"homepage"`
	GitSSHUrl string `json:"git_ssh_url"`
}

type Project struct {
	Name      string `json:"name"`
	WebUrl    string `json:"web_url"`
	GitSSHUrl string `json:"git_ssh_url"`
}

func bindJson(ctx *gin.Context, m interface{}) error {
	err := ctx.BindJSON(m)
	if err != nil {
		ctx.JSON(400, WxResp{ErrCode: 400, ErrMsg: fmt.Sprintf("Parse gitlab requset body error: %s", err)})
		return err
	}
	return nil
}

func buildMsg(content string, markdown bool) string {
	if markdown {
		return fmt.Sprintf(`{"msgtype": "markdown", "markdown":{"content": "%s"}}`, content)
	}
	return fmt.Sprintf(`{"msgtype": "text", "text":{"content": "%s"}}`, content)
}

func TransmitRobot(ctx *gin.Context) {
	key := ctx.GetHeader("X-Gitlab-Token")
	if len(key) == 0 {
		ctx.Render(403, render.Data{ContentType: "application/json", Data: []byte("X-Gitlab-Token is empty")})
		return
	}
	requestUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key)
	var resp *http.Response
	var wxErr error
	var content string
	pushEvent := ctx.GetHeader("X-Gitlab-Event")
	if pushEvent == "Push Hook" {
		pushBody := &PushBody{}
		if err := bindJson(ctx, pushBody); err != nil {
			return
		}
		if len(pushBody.Commits) == 0 && pushBody.After != "0000000000000000000000000000000000000000" {
			ctx.JSON(200, &WxResp{ErrCode: 0, ErrMsg: "no commit"})
			return
		}
		content = "# " + pushBody.Repository.Name + "\n"
		content += "### On branch `" + pushBody.Ref + "`\n"
		for _, v := range pushBody.Commits {
			content += fmt.Sprintf("%s push a commit [%s](%s)  %s", v.Author.Name, strings.ReplaceAll(v.Message, "\n", ""), v.Url, v.TimeStamp) + "\n"
		}
		if pushBody.After == "0000000000000000000000000000000000000000" {
			content += fmt.Sprintf("%s `remove` it", pushBody.UserName)
		}
	} else if pushEvent == "Tag Push Hook" {
		tagPushBody := &TagPushBody{}
		if err := bindJson(ctx, tagPushBody); err != nil {
			return
		}
		content = "# " + tagPushBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s push a tag: [%s](%s)", tagPushBody.UserName, tagPushBody.Ref, tagPushBody.Repository.HomePage+strings.Replace(tagPushBody.Ref, "refs", "", -1))
	} else if pushEvent == "Issue Hook" {
		issueBody := &IssuePushBody{}
		if err := bindJson(ctx, issueBody); err != nil {
			return
		}
		content = "# " + issueBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s %s a issue [%s](%s)", issueBody.User.Name, issueBody.ObjectAttributes.Action, issueBody.ObjectAttributes.Title, issueBody.ObjectAttributes.Url)
	} else if pushEvent == "Note Hook" {
		commentBody := &CommentPushBody{}
		if err := bindJson(ctx, commentBody); err != nil {
			return
		}
		content = "# " + commentBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s leave a comment: %s  %s \n[Detail>>](%s)", commentBody.User.Name, commentBody.ObjectAttributes.Note, commentBody.ObjectAttributes.UpdatedAt, commentBody.ObjectAttributes.Url)
	} else if pushEvent == "Merge Request Hook" {
		mrBody := &MRPushBody{}
		if err := bindJson(ctx, mrBody); err != nil {
			return
		}
		content = "# " + mrBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s `%s` a merge request from `%s` to `%s` \n[Detail>>](%s)", mrBody.User.Name, mrBody.ObjectAttributes.Action, mrBody.ObjectAttributes.SourceBranch, mrBody.ObjectAttributes.TargetBranch, mrBody.ObjectAttributes.Url)
	} else if pushEvent == "Pipeline Hook" {
		pipelineBody := &PipelineBody{}
		if err := bindJson(ctx, pipelineBody); err != nil {
			return
		}
		content = "# " + pipelineBody.Project.Name + "\n"
		branch := "branch"
		if pipelineBody.ObjectAttributes.Tag {
			branch = "tag"
		}
		content += fmt.Sprintf("### Pipeline on %s `%s`\n", branch, pipelineBody.ObjectAttributes.Ref)
		status := ""
		if pipelineBody.ObjectAttributes.Status == "failed" {
			status = "๐"
		} else if pipelineBody.ObjectAttributes.Status == "running" {
			status = "๐"
		} else if pipelineBody.ObjectAttributes.Status == "success" {
			status = "โ"
		} else if pipelineBody.ObjectAttributes.Status == "pending" {
			status = "๐"
		}
		if len(status) == 0 {
			ctx.JSON(200, WxResp{ErrCode: 0, ErrMsg: "unknown status: " + pipelineBody.ObjectAttributes.Status})
			return
		}
		content += "`Status`: " + status + "\n"
		content += fmt.Sprintf("`Start at`: %s\n", pipelineBody.ObjectAttributes.CreatedAt)
		if len(pipelineBody.ObjectAttributes.FinishedAt) > 0 {
			content += fmt.Sprintf("`Finish at`: %s\n", pipelineBody.ObjectAttributes.FinishedAt)
		}
		if pipelineBody.ObjectAttributes.Duration > 0 {
			content += fmt.Sprintf("`Duration`: %ds", pipelineBody.ObjectAttributes.Duration)
		}
	}
	if len(content) == 0 {
		ctx.JSON(200, WxResp{ErrCode: 0, ErrMsg: "no content"})
		return
	}
	content = trans2Emoji(content)
	data := []byte(buildMsg(content, true))
	client := NewClient()
	resp, wxErr = client.Post(requestUrl, "application/json", bytes.NewBuffer(data))
	defer resp.Body.Close()
	if wxErr != nil {
		ctx.JSON(500, WxResp{ErrCode: 500, ErrMsg: fmt.Sprintf("Request wexin robot err: %s ", wxErr)})
		return
	}
	wxResp := &WxResp{}
	json.NewDecoder(resp.Body).Decode(wxResp)
	ctx.JSON(200, wxResp)
}
