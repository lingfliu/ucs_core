package model

import "github.com/lingfliu/ucs_core/utils"

type CtlNodeTemplate struct {
	Id   int64
	Name string

	Template *CtlNode
}

func (nt *CtlNodeTemplate) Validate() bool {

	template := nt.Template

	//非空检测
	if utils.IsEmpty(template.Class) || utils.IsEmpty(template.Name) {
		return false
	}

	return true
}
