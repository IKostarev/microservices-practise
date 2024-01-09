package smtp_client

type Message struct {
	Receivers []string
	Subject   string
	Body      string
}
