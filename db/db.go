package db

import (
	"log"
	"test4/sequence"
)

type SeqDb struct {
}

var group = map[string]string{
	"20210818": "000001",
	"20210817": "000002",
}

// Data 用于获取数据库所有已存数据
func (s *SeqDb) Data() []*sequence.SeqInfo {
	return []*sequence.SeqInfo{}
}

// Save 用于维护对应模块编号数据的相关信息
func (s *SeqDb) Save(model *sequence.SeqInfo) {
	log.Printf("modelInfo %+v", model.Group)
}
