package core

// Context holds shared mutable state for plugins.
type Context struct {
	GlobalState map[string]interface{}
}
