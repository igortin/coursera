package person

var (
	Public  = 1
	private = 1
)

type Person struct {
	Name string
	Id   int
}

//

func (p *Person) ChangeName(val string) {
	p.Name = val
}
