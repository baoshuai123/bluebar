package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//3.ZRevRange 按分数从大到小查询指定数量
	return redisDb.ZRevRange(key, start, end).Result()

}
func GetPostIDsByOrder(p *models.ParamPostList) ([]string, error) {
	//从redis 获取id
	//1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的起始点
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的而数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//使用 pipeline 一次发送多条命令减少 RTT
	pipeline := redisDb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(keyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsByOrder 按社区查询ids
func GetCommunityPostIDsByOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore 把分区的帖子set 与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取数据
	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key 减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if redisDb.Exists(key).Val() < 1 {
		//不存在,需要计算
		pipeline := redisDb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{Aggregate: "MAX"}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}
