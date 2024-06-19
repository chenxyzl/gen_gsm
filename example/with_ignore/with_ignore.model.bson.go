// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.
// gen_tools version: 1.1.4
// generate time: 2024-06-19 16:46:34
package with_ignore

import (
	"github.com/chenxyzl/gsgen/example/common"
	"github.com/chenxyzl/gsgen/gsmodel"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *TestA) MarshalBSON() ([]byte, error) {
	var doc = bson.M{"ig": s.ig}
	return bson.Marshal(doc)
}
func (s *TestA) UnmarshalBSON(data []byte) error {
	doc := struct {
		Ig *common.Common `bson:"ig"`
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetIg(doc.Ig)
	return nil
}
func (s *TestA) BuildBson(m bson.M, preKey string) {
	dirty := s.GetDirty()
	if dirty == 0 {
		return
	}
	if dirty&(1<<0) != 0 {
		if s.ig == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("ig", preKey))
		} else {
			s.ig.BuildBson(m, gsmodel.MakeBsonKey("ig", preKey))
		}
	}
	return
}
func (s *TestA) CleanDirty() {
	s.DirtyModel.CleanDirty()
	if s.ig != nil {
		s.ig.CleanDirty()
	}
}
