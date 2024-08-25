package model

var PlatformIcons = map[string]string{
	"zenn":          "internal/server/assets/zenn.png",
	"qiita":         "internal/server/assets/qiita.png",
	"twitter":       "internal/server/assets/twitter.png",
	"linkedin":      "internal/server/assets/linkedin.png",
	"stackoverflow": "internal/server/assets/stackoverflow.png",
	"atcoder":       "internal/server/assets/atcoder.png",
	"note":          "internal/server/assets/note.png",
	"youtube":       "internal/server/assets/youtube.png",
	"instagram":     "internal/server/assets/instagram.png",
	"facebook":      "internal/server/assets/facebook.png",
}

var PlatformURLs = map[string]string{
	"zenn":          "https://zenn.dev/",
	"qiita":         "https://qiita.com/",
	"twitter":       "https://x.com/",
	"linkedin":      "https://linkedin.com/in/",
	"stackoverflow": "https://stackoverflow.com/users/",
	"atcoder":       "https://atcoder.jp/users/",
	"note":          "https://note.com/",
	"youtube":       "https://youtube.com/",
	"instagram":     "https://instagram.com/",
	"facebook":      "https://facebook.com/",
}

var PlatformColors = map[string]string{
	"zenn":          "#3EA8FF",
	"qiita":         "#55C500",
	"twitter":       "#FFFFFF",
	"linkedin":      "#0A66C2",
	"stackoverflow": "#F48024",
	"atcoder":       "#000000",
	"note":          "#00A7AF",
	"youtube":       "#FF0000",
	"instagram":     "#E4405F",
	"facebook":      "#1877F2",
}

var PlatformBgColors = map[string]string{
	"zenn":          "#F1F5F9",
	"qiita":         "#F5F6F6",
	"twitter":       "#000000",
	"linkedin":      "#F4F2EE",
	"stackoverflow": "#FFFFFB",
	"atcoder":       "#EBEBEB",
	"note":          "#F5F5F5",
	"youtube":       "#FFFFFF",
	"instagram":     "#F5F5F5",
	"facebook":      "#F5F5F5",
}

var PlatformFontColors = map[string]string{
	"zenn":          "#000000",
	"qiita":         "#000000",
	"twitter":       "#FFFFFF",
	"linkedin":      "#000000",
	"stackoverflow": "#000000",
	"atcoder":       "#000000",
	"note":          "#000000",
	"youtube":       "#000000",
	"instagram":     "#000000",
	"facebook":      "#000000",
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
	Ranking        int    // AtCoder用のフィールド
	Rating         int    // AtCoder用のフィールド
	RatedMatches   int    // AtCoder用のフィールド
	CustomURL      string // Youtube用のフィールド
	TotalVideos    int    // Youtube用のフィールド
	TotalViewCount int    // Youtube用のフィールド
	Experience     string // Linkedin用のフィールド
	Education      string // Linkedin用のフィールド
	Location       string // Linkedin用のフィールド
}
