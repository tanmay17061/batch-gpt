package config

type ServingMode interface {
    IsAsync() bool
    IsCache() bool
    GetMode() string
}

type servingMode struct {
    mode string
}

func NewServingMode(mode string) ServingMode {
    if mode == "" {
        mode = "sync" // Default to synchronous mode
    }
    return &servingMode{mode: mode}
}

func (sm *servingMode) IsAsync() bool {
    return sm.mode == "async"
}

func (sm *servingMode) IsCache() bool {
    return sm.mode == "cache"
}

func (sm *servingMode) GetMode() string {
    return sm.mode
}
