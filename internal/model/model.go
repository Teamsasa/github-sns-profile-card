package model

var PlatformIcons = map[string]string{
	"zenn":          "/assets/zenn.png",
	"qiita":         "/assets/qiita.png",
	"twitter":       "/assets/twitter.png",
	"linkedin":      "/assets/linkedin.png",
	"stackoverflow": "/assets/stackoverflow.png",
	"atcoder":       "/assets/atcoder.png",
}

var PlatformURLs = map[string]string{
	"zenn":          "https://zenn.dev/",
	"qiita":         "https://qiita.com/",
	"twitter":       "https://x.com/",
	"linkedin":      "https://linkedin.com/in/",
	"stackoverflow": "https://stackoverflow.com/users/",
	"atcoder":       "https://atcoder.jp/users/",
}

var PlatformColors = map[string]string{
	"zenn":          "#3EA8FF",
	"qiita":         "#55C500",
	"twitter":       "#FFFFFF",
	"linkedin":      "#0A66C2",
	"stackoverflow": "#F48024",
	"atcoder":       "#000000",
}

var PlatformBgColors = map[string]string{
	"zenn":          "#F1F5F9",
	"qiita":         "#F5F6F6",
	"twitter":       "#000000",
	"linkedin":      "#F4F2EE",
	"stackoverflow": "#FFFFFB",
	"atcoder":       "#EBEBEB",
}

var PlatformFontColors = map[string]string{
	"zenn":          "#000000",
	"qiita":         "#000000",
	"twitter":       "#FFFFFF",
	"linkedin":      "#000000",
	"stackoverflow": "#000000",
	"atcoder":       "#000000",
}

type PlatformUserInfo struct {
	FollowersCount int
	FollowingCount int
	ArticlesCount  int
	LikeCount      int    // Zenn用のフィールド
	UserName       string // Stackoverflow用のフィールド
	Reputation     int    // StackOverflow用のフィールド
	AnswerCount    int    // StackOverflow用のフィールド
	QuestionCount  int    // StackOverflow用のフィールド
	Rating         int    // AtCoder用のフィールド
}
