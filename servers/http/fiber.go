package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net"
)

type FiberApp struct {
	app *fiber.App
}

func (s *FiberApp) Get(path string, handlers ...Handler) Router {
	s.app.Get(path, FiberWrapHandlers(handlers...)...)
	return s
}

func (s *FiberApp) Head(path string, handlers ...Handler) Router {
	s.app.Head(path, FiberWrapHandlers(handlers...)...)
	return s
}

func (s *FiberApp) Post(path string, handlers ...Handler) Router {
	s.app.Post(path, FiberWrapHandlers(handlers...)...)
	return s
}

func (s *FiberApp) Options(path string, handlers ...Handler) Router {
	s.app.Options(path, FiberWrapHandlers(handlers...)...)
	return s
}

func (s *FiberApp) Delete(path string, handlers ...Handler) Router {
	s.app.Delete(path, FiberWrapHandlers(handlers...)...)
	return s
}

func (s *FiberApp) Use(args ...interface{}) Router {
	s.app.Use(args)
	return s
}

func (s *FiberApp) Group(prefix string, handlers ...Handler) Router {
	gr := s.app.Group(prefix, FiberWrapHandlers(handlers...)...)
	return NewFiberGroup(gr)
}

func (s *FiberApp) Listener(ln net.Listener) error {
	return s.app.Listener(ln)
}

func (s *FiberApp) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *FiberApp) Shutdown() error {
	return s.app.Shutdown()
}

func FiberWrapHandlers(handlers ...Handler) []fiber.Handler {
	var fiberHandlers []fiber.Handler
	for _, handler := range handlers {
		handler := handler // copy
		fiberHandlers = append(fiberHandlers, func(ctx *fiber.Ctx) error {
			return handler(newFiberContext(ctx))
		})
	}

	return fiberHandlers
}

// NewFiberServer - return wrapper of Fiber App
func NewFiberServer(f *fiber.App) Server {
	return &FiberApp{app: f}
}

// FiberContext - wrapper on context fiber lib
type FiberContext struct {
	context *fiber.Ctx
	request Request
}

func (f *FiberContext) IP() string {
	return f.context.IP()
}

func (f *FiberContext) Hostname() string {
	return f.context.Hostname()
}

func (f *FiberContext) Query(key string, defaultValue ...string) string {
	return f.context.Query(key, defaultValue...)
}

func (f *FiberContext) Set(key string, val string) {
	f.context.Set(key, val)
}

func (f *FiberContext) Append(field string, values ...string) {
	f.context.Append(field, values...)
}

func (f *FiberContext) Write(p []byte) (int, error) {
	return f.context.Write(p)
}

func (f *FiberContext) Status(status int) Context {
	f.context.Status(status)
	return f
}

func (f *FiberContext) GetReqHeaders() map[string]string {
	return f.context.GetReqHeaders()
}

func (f *FiberContext) Request() Request {
	return f.request
}

func (f *FiberContext) Writef(s string, a ...interface{}) (int, error) {
	return f.context.Writef(s, a...)
}

func (f *FiberContext) WriteString(s string) (int, error) {
	return f.context.WriteString(s)
}

func (f *FiberContext) BodyParser(out interface{}) error {
	return f.context.BodyParser(out)
}

func (f *FiberContext) Next() error {
	return f.context.Next()
}

func (f *FiberContext) Redirect(location string, status int) error {
	return f.context.Redirect(location, status)
}

func (f *FiberContext) Locals(key string, value ...interface{}) (val interface{}) {
	return f.context.Locals(key, value...)
}

func (f *FiberContext) Get(key string, defaultValue ...string) string {
	return f.context.Get(key, defaultValue...)
}

func (f *FiberContext) Method(override ...string) string {
	return f.context.Method(override...)
}

func (f *FiberContext) Params(key string, defaultValue ...string) string {
	return f.context.Params(key, defaultValue...)
}

func newFiberContext(ctx *fiber.Ctx) *FiberContext {
	return &FiberContext{
		context: ctx,
		request: newFiberRequest(ctx.Request()),
	}
}

// FiberRequest - wrapper on fiber fasthttp request
type FiberRequest struct {
	request *fasthttp.Request
}

func (f *FiberRequest) GetContentLength() int {
	return f.request.Header.ContentLength()
}

func (f *FiberRequest) Body() []byte {
	return f.request.Body()
}

func (f *FiberRequest) RequestURI() string {
	return string(f.request.RequestURI())
}

func newFiberRequest(r *fasthttp.Request) Request {
	return &FiberRequest{request: r}
}

type FiberGroup struct {
	gr fiber.Router
}

func (fg *FiberGroup) Get(path string, handlers ...Handler) Router {
	fg.gr.Get(path, FiberWrapHandlers(handlers...)...)
	return fg
}

func (fg *FiberGroup) Head(path string, handlers ...Handler) Router {
	fg.gr.Head(path, FiberWrapHandlers(handlers...)...)
	return fg
}

func (fg *FiberGroup) Post(path string, handlers ...Handler) Router {
	fg.gr.Post(path, FiberWrapHandlers(handlers...)...)
	return fg
}

func (fg *FiberGroup) Options(path string, handlers ...Handler) Router {
	fg.gr.Options(path, FiberWrapHandlers(handlers...)...)
	return fg
}

func (fg *FiberGroup) Delete(path string, handlers ...Handler) Router {
	fg.gr.Delete(path, FiberWrapHandlers(handlers...)...)
	return fg
}

func (fg *FiberGroup) Use(args ...interface{}) Router {
	fg.gr.Use(args)
	return fg
}

func (fg *FiberGroup) Group(prefix string, handlers ...Handler) Router {
	gr := fg.gr.Group(prefix, FiberWrapHandlers(handlers...)...)
	return NewFiberGroup(gr)
}

func NewFiberGroup(gr fiber.Router) *FiberGroup {
	return &FiberGroup{gr: gr}
}
