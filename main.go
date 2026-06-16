package main

import (
"context"
"fmt"
"time"

"github.com/redis/go-redis/v9"
)

func main() {
ctx := context.Background()

// 连接到 Docker Compose 中启动的真实 Redis
client := redis.NewClient(&redis.Options{
Addr: "redis:6379", // 走 docker 容器网络
})
defer client.Close()

fmt.Println("清理历史慢查询记录...")
client.Do(ctx, "SLOWLOG", "RESET")

fmt.Println("正在制造 2 个大 Key (包含 1,000,000 个元素)，这需要一点时间...")
createBigKey(ctx, client, "bigkey_for_del", 1000000)
createBigKey(ctx, client, "bigkey_for_unlink", 1000000)

fmt.Println("\n大 Key 制造完毕。开始删除测试：")

// 1. 测试 DEL
fmt.Println(">>> 1. 正在使用 DEL 指令删除 bigkey_for_del ...")
startDel := time.Now()
client.Del(ctx, "bigkey_for_del")
fmt.Printf(">>> DEL 客户端耗时结束，经过了 %v\n", time.Since(startDel))

// 2. 测试 UNLINK
fmt.Println(">>> 2. 正在使用 UNLINK 指令删除 bigkey_for_unlink ...")
startUnlink := time.Now()
client.Unlink(ctx, "bigkey_for_unlink")
fmt.Printf(">>> UNLINK 客户端耗时结束，经过了 %v\n", time.Since(startUnlink))

// 获取慢查询日志
fmt.Println("\n=== 查询 Redis 慢查询日志 ===")
slowlogs, _ := client.Do(ctx, "SLOWLOG", "GET", 5).Slice()

if len(slowlogs) == 0 {
fmt.Println("未发现慢查询日志。")
return
}

for i, logItem := range slowlogs {
log := logItem.([]interface{})
duration := log[2]
commands := log[3].([]interface{})

fmt.Printf("%d) 执行时长: %v 微秒 | 执行命令: ", i+1, duration)
for _, cmd := range commands {
fmt.Printf("%s ", cmd)
}
fmt.Println()
}
}

func createBigKey(ctx context.Context, client *redis.Client, keyName string, size int) {
pipe := client.Pipeline()
for i := 0; i < size; i++ {
pipe.SAdd(ctx, keyName, fmt.Sprintf("val_%d", i))
// 每 10 万条提交一次，防止占用过多内存
if i > 0 && i%100000 == 0 {
pipe.Exec(ctx)
}
}
pipe.Exec(ctx)
}
