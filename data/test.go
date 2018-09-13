package data

type Test struct {
	ID 	 int 
	Desc string
}

func (t Test) GetID() interface{} {
	return t.ID
}