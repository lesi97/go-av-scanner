package store

type Context struct {
	Foo string
}

type apiContextKeyType string

const ContextKey apiContextKeyType = "api"