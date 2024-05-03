# Golang shared lock pool
golang library for lock pool implementation by using mutex
```
go get github.com/radityacandra/lock-pool
```

## Usage
- import it

```
import lp "github.com/radityacandra/lock-pool"
```

- initiate
```
lock_pool := lp.NewLockPool()
```

- use
```
ctx := context.Background()
lock_pool.Acquire(ctx, "some-lock")
// do some task here
lock_pool.Release(ctx, "some-lock")
```
