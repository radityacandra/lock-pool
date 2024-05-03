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
context.Background()
lp.Acquire(ctx, "some-lock")
// do some task here
lp.Release(ctx, "some-lock")
```
