package resources

type DataCenter struct {
	Name        string
	Id          string
	Region      string
	Description string
}

func (d *DataCenter) GetName() string {
	return d.Name
}

func (d *DataCenter) GetId() string {
	return d.Id
}

func (d *DataCenter) GetRegion() string {
	return d.Region
}

func (d *DataCenter) GetDescription() string {
	return d.Description
}
