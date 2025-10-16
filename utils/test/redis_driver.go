package test

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient is a mock implementation of redis.Client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) AddHook(hook redis.Hook) {
	m.Called(hook)
}

func (m *MockRedisClient) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	args := m.Called(ctx, fn, keys)
	return args.Error(0)
}

func (m *MockRedisClient) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	mockArgs := m.Called(ctx, args)
	return mockArgs.Get(0).(*redis.Cmd)
}

func (m *MockRedisClient) Process(ctx context.Context, cmd redis.Cmder) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockRedisClient) ProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	args := m.Called(ctx, cmds)
	return args.Error(0)
}

func (m *MockRedisClient) ProcessMulti(ctx context.Context, cmds []redis.Cmder) error {
	args := m.Called(ctx, cmds)
	return args.Error(0)
}

func (m *MockRedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	args := m.Called(ctx, channels)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisClient) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	args := m.Called(ctx, channels)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisClient) SSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	args := m.Called(ctx, channels)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) PoolStats() *redis.PoolStats {
	args := m.Called()
	return args.Get(0).(*redis.PoolStats)
}

func (m *MockRedisClient) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	args := m.Called(ctx, fn)
	return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockRedisClient) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	args := m.Called(ctx, fn)
	return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockRedisClient) TxPipeline() redis.Pipeliner {
	args := m.Called()
	return args.Get(0).(redis.Pipeliner)
}

func (m *MockRedisClient) Pipeline() redis.Pipeliner {
	args := m.Called()
	return args.Get(0).(redis.Pipeliner)
}

func (m *MockRedisClient) ClientGetName(ctx context.Context) *redis.StringCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	args := m.Called(ctx, message)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Quit(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Dump(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	args := m.Called(ctx, key, tm)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	args := m.Called(ctx, pattern)
	return args.Get(0).(*redis.StringSliceCmd)
}

func (m *MockRedisClient) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, host, port, key, db, timeout)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Move(ctx context.Context, key string, db int) *redis.BoolCmd {
	args := m.Called(ctx, key, db)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.DurationCmd)
}

func (m *MockRedisClient) Persist(ctx context.Context, key string) *redis.BoolCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	args := m.Called(ctx, key, tm)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.DurationCmd)
}

func (m *MockRedisClient) RandomKey(ctx context.Context) *redis.StringCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	args := m.Called(ctx, key, newkey)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	args := m.Called(ctx, key, newkey)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	args := m.Called(ctx, key, ttl, value)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	args := m.Called(ctx, key, ttl, value)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	args := m.Called(ctx, key, sort)
	return args.Get(0).(*redis.StringSliceCmd)
}

func (m *MockRedisClient) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	args := m.Called(ctx, key, store, sort)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	args := m.Called(ctx, key, sort)
	return args.Get(0).(*redis.SliceCmd)
}

func (m *MockRedisClient) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.DurationCmd)
}

func (m *MockRedisClient) Type(ctx context.Context, key string) *redis.StatusCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	args := m.Called(ctx, cursor, match, count)
	return args.Get(0).(*redis.ScanCmd)
}

func (m *MockRedisClient) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	args := m.Called(ctx, key, cursor, match, count)
	return args.Get(0).(*redis.ScanCmd)
}

func (m *MockRedisClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	args := m.Called(ctx, key, cursor, match, count)
	return args.Get(0).(*redis.ScanCmd)
}

func (m *MockRedisClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	args := m.Called(ctx, key, cursor, match, count)
	return args.Get(0).(*redis.ScanCmd)
}

func (m *MockRedisClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	args := m.Called(ctx, key, bitCount)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, destKey, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, destKey, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, destKey, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) BitOpNot(ctx context.Context, destKey, key string) *redis.IntCmd {
	args := m.Called(ctx, destKey, key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	args := m.Called(ctx, key, bit, pos)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	args := m.Called(ctx, key, decrement)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	args := m.Called(ctx, key, offset)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	args := m.Called(ctx, key, start, end)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*redis.FloatCmd)
}

func (m *MockRedisClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.SliceCmd)
}

func (m *MockRedisClient) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	args := m.Called(ctx, values)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	args := m.Called(ctx, values)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	args := m.Called(ctx, key, offset, value)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	args := m.Called(ctx, key, offset, value)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

// MockRedisCmd is a mock implementation of redis.Cmd
type MockRedisCmd struct {
	mock.Mock
}

func (m *MockRedisCmd) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRedisCmd) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRedisCmd) FullName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRedisCmd) Args() []interface{} {
	args := m.Called()
	return args.Get(0).([]interface{})
}

func (m *MockRedisCmd) Val() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockRedisCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisCmd) readReply(rd interface{}) error {
	args := m.Called(rd)
	return args.Error(0)
}

func (m *MockRedisCmd) readTimeout() *time.Duration {
	args := m.Called()
	return args.Get(0).(*time.Duration)
}

func (m *MockRedisCmd) SetFirstKeyPos(keyPos int8) {
	m.Called(keyPos)
}

func (m *MockRedisCmd) SetReadTimeout(timeout time.Duration) {
	m.Called(timeout)
}

func (m *MockRedisCmd) SetWriteTimeout(timeout time.Duration) {
	m.Called(timeout)
}

// MockRedisStatusCmd is a mock implementation of redis.StatusCmd
type MockRedisStatusCmd struct {
	MockRedisCmd
}

func (m *MockRedisStatusCmd) Val() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRedisStatusCmd) Result() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// MockRedisStringCmd is a mock implementation of redis.StringCmd
type MockRedisStringCmd struct {
	MockRedisCmd
}

func (m *MockRedisStringCmd) Val() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRedisStringCmd) Result() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// MockRedisIntCmd is a mock implementation of redis.IntCmd
type MockRedisIntCmd struct {
	MockRedisCmd
}

func (m *MockRedisIntCmd) Val() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockRedisIntCmd) Result() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockRedisBoolCmd is a mock implementation of redis.BoolCmd
type MockRedisBoolCmd struct {
	MockRedisCmd
}

func (m *MockRedisBoolCmd) Val() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRedisBoolCmd) Result() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

// MockRedisFloatCmd is a mock implementation of redis.FloatCmd
type MockRedisFloatCmd struct {
	MockRedisCmd
}

func (m *MockRedisFloatCmd) Val() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *MockRedisFloatCmd) Result() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

// MockRedisStringSliceCmd is a mock implementation of redis.StringSliceCmd
type MockRedisStringSliceCmd struct {
	MockRedisCmd
}

func (m *MockRedisStringSliceCmd) Val() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockRedisStringSliceCmd) Result() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// MockRedisSliceCmd is a mock implementation of redis.SliceCmd
type MockRedisSliceCmd struct {
	MockRedisCmd
}

func (m *MockRedisSliceCmd) Val() []interface{} {
	args := m.Called()
	return args.Get(0).([]interface{})
}

func (m *MockRedisSliceCmd) Result() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockRedisDurationCmd is a mock implementation of redis.DurationCmd
type MockRedisDurationCmd struct {
	MockRedisCmd
}

func (m *MockRedisDurationCmd) Val() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockRedisDurationCmd) Result() (time.Duration, error) {
	args := m.Called()
	return args.Get(0).(time.Duration), args.Error(1)
}

// MockRedisScanCmd is a mock implementation of redis.ScanCmd
type MockRedisScanCmd struct {
	MockRedisCmd
}

func (m *MockRedisScanCmd) Val() (keys []string, cursor uint64) {
	args := m.Called()
	return args.Get(0).([]string), args.Get(1).(uint64)
}

func (m *MockRedisScanCmd) Result() (keys []string, cursor uint64, err error) {
	args := m.Called()
	return args.Get(0).([]string), args.Get(1).(uint64), args.Error(2)
}

// MockRedisPubSub is a mock implementation of redis.PubSub
type MockRedisPubSub struct {
	mock.Mock
}

func (m *MockRedisPubSub) Subscribe(ctx context.Context, channels ...string) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRedisPubSub) PSubscribe(ctx context.Context, channels ...string) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRedisPubSub) SSubscribe(ctx context.Context, channels ...string) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRedisPubSub) Unsubscribe(ctx context.Context, channels ...string) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRedisPubSub) PUnsubscribe(ctx context.Context, channels ...string) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRedisPubSub) SUnsubscribe(ctx context.Context, channels ...string) error {
	args := m.Called(ctx, channels)
	return args.Error(0)
}

func (m *MockRedisPubSub) Ping(ctx context.Context, payload string) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockRedisPubSub) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisPubSub) Receive(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func (m *MockRedisPubSub) ReceiveMessage(ctx context.Context) (*redis.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(*redis.Message), args.Error(1)
}

func (m *MockRedisPubSub) ReceiveTimeout(ctx context.Context, timeout time.Duration) (interface{}, error) {
	args := m.Called(ctx, timeout)
	return args.Get(0), args.Error(1)
}

func (m *MockRedisPubSub) Channel() <-chan *redis.Message {
	args := m.Called()
	return args.Get(0).(<-chan *redis.Message)
}

func (m *MockRedisPubSub) ChannelSize(size int) <-chan *redis.Message {
	args := m.Called(size)
	return args.Get(0).(<-chan *redis.Message)
}

func (m *MockRedisPubSub) ChannelWithSubscriptions(ctx context.Context, bufferSize int) <-chan interface{} {
	args := m.Called(ctx, bufferSize)
	return args.Get(0).(<-chan interface{})
}

func (m *MockRedisPubSub) WithChannels(channels ...string) *redis.PubSub {
	args := m.Called(channels)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisPubSub) WithChannelSize(size int) *redis.PubSub {
	args := m.Called(size)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisPubSub) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRedisPubSub) SetVal(val interface{}) {
	m.Called(val)
}

func (m *MockRedisPubSub) SetErr(e error) {
	m.Called(e)
}

// MockRedisPipeliner is a mock implementation of redis.Pipeliner
type MockRedisPipeliner struct {
	mock.Mock
}

func (m *MockRedisPipeliner) AddHook(hook redis.Hook) {
	m.Called(hook)
}

func (m *MockRedisPipeliner) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	mockArgs := m.Called(ctx, args)
	return mockArgs.Get(0).(*redis.Cmd)
}

func (m *MockRedisPipeliner) Process(ctx context.Context, cmd redis.Cmder) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockRedisPipeliner) ProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	args := m.Called(ctx, cmds)
	return args.Error(0)
}

func (m *MockRedisPipeliner) ProcessMulti(ctx context.Context, cmds []redis.Cmder) error {
	args := m.Called(ctx, cmds)
	return args.Error(0)
}

func (m *MockRedisPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
	args := m.Called(ctx)
	return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockRedisPipeliner) ExecVal(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	args := m.Called(ctx, fn)
	return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockRedisPipeliner) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisPipeliner) Discard() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisPipeliner) Len() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockRedisPipeliner) SetVal(val interface{}) {
	m.Called(val)
}

func (m *MockRedisPipeliner) SetErr(e error) {
	m.Called(e)
}
