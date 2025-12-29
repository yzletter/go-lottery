# go-lottery

## 压力测试

QPS 2200, avg time 59ms

```shell
yzletter@yangzhileideMacBook-Pro go-lottery % go test -v ./main_test.go -run=^TestLottery$ -count=1
=== RUN   TestLottery
QPS 2200, avg time 59ms
27      100
25      1000
20      1000
24      400
26      1000
1       1000
21      1000
28      500
22      200
23      300
共计6500件商品
--- PASS: TestLottery (3.95s)
PASS
ok      command-line-arguments  4.810s
```

