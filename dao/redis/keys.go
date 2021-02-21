package redis

// redis key

// redis key 注意使用命名空间的方式,方便查询和拆分
const (
	Prefix                 = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //zset;帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //zset;帖子及投票分数
	keyPostVotedZSetPrefix = "post:voted:" //zset;记录用户及投票类型;参数是post
	KeyCommunitySetPF      = "community:"  //set;保存每个分区下帖子id
)

//给redis key 加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
