package args

import "github.com/go-redis/redis/v8"

func (ra *RArgs) BuildZArgs(score float64, member interface{}) *redis.ZAddArgs {
	rz := &redis.Z{Score: score, Member: member}
	za := &redis.ZAddArgs{
		XX: false, //只更新已经存在的元素。不要添加新元素。
		NX: true,  //只添加新元素。不要更新已经存在的元素。
		LT: false, //：仅当新分数低于当前分数时才更新现有元素。此标志不会阻止添加新元素。
		GT: true,  //：仅当新分数大于当前分数时才更新现有元素。此标志不会阻止添加新元素。
		Ch: true,  //：将返回值从添加的新元素数修改为改变的元素总数（CH 是changed的​​缩写）。更改的元素是添加的新元素和已更新分数的元素。因此，命令行中指定的与过去得分相同的元素不计算在内。注意：通常返回值ZADD只计算添加的新元素的数量。
		// INCR: false, //：指定此选项时的ZADD行为类似于ZINCRBY。在此模式下只能指定一个分数元素对。
		Members: []redis.Z{*rz},
	}
	return za
}
