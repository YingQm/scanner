//提供限流功能：在指定的时间周期内每个IP最多只能访问指定次数
package util

import (
	"errors"
	lru "github.com/hashicorp/golang-lru"
	"time"
)

type AccessInfo struct {
	IP             string
	PrevCount      int64 //前一个周期访问次数
	CurStartTime   int64 //当前周期开始时间unix时间戳
	CurLastAccTime int64 //当前周期最新访问时间unix时间戳
	Count          int64 //当前访问次数
}

type RateLimiter struct {
	TimeInterval int64 //时间间隔（单位秒）
	MaxCount     int64 //最大访问次数
	Cache        *lru.Cache
}

func NewRateLimiter(interval int64, count int64) (*RateLimiter, error) {
	var (
		cache   *lru.Cache
		limiter *RateLimiter
		err     error
	)
	if interval <= 0 || count <= 0 {
		return nil, errors.New("invalinfo argument")
	}

	if cache, err = lru.New(10000); err != nil {
		return nil, err
	}

	limiter = &RateLimiter{
		MaxCount:     count,
		TimeInterval: interval,
		Cache:        cache,
	}

	return limiter, nil
}

//根据标识符判断是否允许访问
func (limit *RateLimiter) Allow(ip string) bool {
	var now = time.Now().Unix()
	ainfo, ok := limit.Cache.Get(ip)
	if ok {
		info := ainfo.(AccessInfo)

		//距离上次访问的时间间隔
		interval := now - info.CurLastAccTime
		if interval < 0 {
			return false
		}

		if interval >= limit.TimeInterval {
			//距上次访问已经超过一个周期，重置上个周期的访问次数为0
			info.PrevCount = 0
			info.CurStartTime = now
			info.CurLastAccTime = now
			info.Count = 1
			limit.Cache.Add(ip, info)
			return true
		}

		//判断距离开始时间是否在一个时间周期内
		timespan := now - info.CurStartTime
		if timespan <= limit.TimeInterval {
			if info.Count+1 > limit.MaxCount {
				return false
			}
			//根据上个周期的访问次数和当前已访问次数估算
			//从当前时间往前推一个时间周期的访问次数
			estimateCount := info.PrevCount*(limit.TimeInterval-timespan)/limit.TimeInterval + info.Count
			if estimateCount > limit.MaxCount {
				return false
			}
			info.Count++
			info.CurLastAccTime = now
			limit.Cache.Add(ip, info)
			return true
		} else {
			//当前处于新的计数周期
			info.PrevCount = info.Count
			info.CurStartTime = now
			info.CurLastAccTime = now
			info.Count = 1
			limit.Cache.Add(ip, info)
			return true
		}
	}
	//缓存中不存在，新增
	info := AccessInfo{}
	info.IP = ip
	info.PrevCount = 0
	info.CurStartTime = time.Now().Unix()
	info.CurLastAccTime = time.Now().Unix()
	info.Count = 1
	limit.Cache.Add(ip, info)
	return true
}
