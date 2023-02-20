package ecs

type IComponent interface {
	GetId() string
	GetData() interface{}
	SetData(v interface{})
	GetStructure() *Component
}

type Component struct {
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}

func CreateComponent(id string, data interface{}) *Component {
	return &Component{Id: id, Data: data}
}

func (p *Component) GetId() string {
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
