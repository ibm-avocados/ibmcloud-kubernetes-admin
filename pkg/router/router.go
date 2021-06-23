package router

type Rest interface {
	Serve(string) error
}
