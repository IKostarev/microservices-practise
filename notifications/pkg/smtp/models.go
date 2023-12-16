package smtp

type Message struct {
	Receivers []string
	Subject   string
	Body      string
}
