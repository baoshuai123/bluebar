package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePreVote     = 432 //每一票分数
)

var (
	ErrVoteTimeExpire = errors.New("超出投票时间")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := redisDb.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis 取有序集合    postID帖子的发布时间
	postTime := redisDb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2 和 3 需要进行事务操作
	//2.更新分数
	//先查当前用户给当前帖子的投票记录
	ov := redisDb.ZScore(getRedisKey(keyPostVotedZSetPrefix+postID), userID).Val()
	//更新:如果这一次投票和之前的保存的值一致,就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票差值
	pipeline := redisDb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePreVote, postID)
	//3.记录用户为帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(keyPostVotedZSetPrefix+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(keyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value, //赞成还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
