package typedast

type Stm interface {
	Node
	stmNode()
}

type BaseStm struct {
	line int
	col  int
	text string
}

func (s BaseStm) Line() int    { return s.line }
func (s BaseStm) Col() int     { return s.col }
func (s BaseStm) Text() string { return s.text }

type ExpStm struct {
	Exp Exp

	BaseStm
}

func (*ExpStm) stmNode() {}

func NewExpStm(
	exp Exp,
	line int,
	col int,
	text string,
) ExpStm {
	return ExpStm{
		Exp: exp,
		BaseStm: BaseStm{
			line: line,
			col:  col,
			text: text,
		},
	}
}

// check that ExpStm implements Stm
var _ Stm = (*ExpStm)(nil)
