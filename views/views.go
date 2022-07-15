package views

type Engine interface {
	Render(string, map[string]string) (string, error)
}
