package args

var (
	RArgsApp = NewArgs()
)

type RArgs struct{}

func NewArgs() *RArgs {
	return &RArgs{}
}
