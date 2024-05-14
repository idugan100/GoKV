# GoKV
![Tests/Linting/Formatting/Security](https://github.com/idugan100/GoKV/actions/workflows/main.yml/badge.svg)


This is a simple in-memory key-value store written in Go. All it's commands are redis compatible (can be used via Redis CLI or any redis library), but the full redis 7.0 spec has not been implemented. This is not production ready. It was built as a learning exercise to better understand redis internals. I used a number of resources to learn how redis works in order to built GoKV that can be found [here](./Resources.MD).

# Suported Commands
- [PING](https://redis.io/docs/latest/commands/ping/)
- [SET](https://redis.io/docs/latest/commands/set/)
- [DEL](https://redis.io/docs/latest/commands/del/)
- [RANDOMKEY](https://redis.io/docs/latest/commands/randkey/)
- [EXISTS](https://redis.io/docs/latest/commands/exists/)
- [STRLEN](https://redis.io/docs/latest/commands/strlen/)
- [LOLWUT](https://redis.io/docs/latest/commands/lolwut/)
- [FLUSHALL](https://redis.io/docs/latest/commands/flushall/)
- [GETSET](https://redis.io/docs/latest/commands/getset/)