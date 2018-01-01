package gologo

type Lifecycle []func()

func New() *Lifecycle {
    s := make(Lifecycle, 0)
    return &s
}

func (lifecycle *Lifecycle) Cleanup() {
    for _, callback := range *lifecycle {
        callback()
    }
}

func (lifecycle *Lifecycle) RegisterCleanup(callback func()) {
    s := append(*lifecycle, callback)
    lifecycle = &s
}