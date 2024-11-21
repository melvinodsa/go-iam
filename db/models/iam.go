package models

type iam struct{}

func (i iam) DbName() string {
	return "iam"
}
