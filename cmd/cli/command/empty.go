package command

type Empty struct{}

func NewEmpty() (*Empty, error) {
	return &Empty{}, nil
}

func (cmd *Empty) Do() string {
	return "did nothing\n"
}
