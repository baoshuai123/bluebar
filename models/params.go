package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//定义请求的参数结构体
//注册参数
type ParamSignUP struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

//登录参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//投票参数
type ParamVoteData struct {
	// UserID 从当前请求获取
	PostID    string `json:"post_id"  binding:"required"`             //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` //赞成票(1)反对票(-1)取消投票(0)
}

// ParamPostList 获取帖子列表 query string参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

// ParamCommunityPostList 按社区获取帖子列表 query string参数
