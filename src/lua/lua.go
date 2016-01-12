package lua

import(
	"github.com/garyburd/redigo/redis"
	scripts "lua/scripts"
)

type Lua struct {
	redisPool *redis.Pool
	consumeUsage *redis.Script
}

func NewClient(rp *redis.Pool) *Lua {
	return &Lua{
		redisPool: rp,
		consumeUsage: redis.NewScript(0, scripts.CONSUME_USAGE),
		
	}
}

func (l *Lua) ConsumeTokenUsage(token string) (usesLeft int, err error) {
	r := l.redisPool.Get()
	defer r.Close()

	usesLeft, err = redis.Int(l.consumeUsage.Do(r, token))
	return
}
