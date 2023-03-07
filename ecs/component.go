package ecs

type IComponent interface {
	GetId() Identifier
	GetData() interface{}
	SetData(v interface{})
	GetStructure() *Component
}

type Component struct {
	Id   Identifier  `json:"id"`
	Data interface{} `json:"data"`
}

func CreateComponent(id Identifier, data interface{}) *Component {
	return &Component{Id: id, Data: data}
}

func (p *Component) GetId() Identifier {
	return p.Id
}

func (p *Component) GetData() interface{} {
	return p.Data
}

func (p *Component) SetData(v interface{}) {
	p.Data = v
}

func (p *Component) GetStructure() *Component {
	return p
}
