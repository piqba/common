# Tools module

This module, pkg or folder contain customs implementations for commons uses cases like a data structure `Set` or `time checking` scenarios 


```go
// Set - our representation of set data structure
type Set struct {
    Elements map[Any]struct{}
}

// TimeChecker ...
type TimeChecker struct {
    Start  time.Time
    End    time.Time
    Check  time.Time
    Tz     string
    Layout string
}
```