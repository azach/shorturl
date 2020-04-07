package cache

type Cache interface {
	Get(key string) (value string, exists bool)
	Set(key string, value string) error
}
