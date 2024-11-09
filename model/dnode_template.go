package model

import "github.com/lingfliu/ucs_core/utils"

type DNodeTemplate struct {
	Id   int64
	Name string

	Template *DNode
}

func (dnt *DNodeTemplate) Validate() bool {

	template := dnt.Template

	//非空检测
	if utils.IsEmpty(template.Class) || utils.IsEmpty(template.Name) {
		return false
	}

	return true
}
