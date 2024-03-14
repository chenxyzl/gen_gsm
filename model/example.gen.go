// Code generated by genmod; DO NOT EDIT.
// genmod version: 0.0.1
// generate time: 2024-03-14 16:22:39
// src code version: 
// src code commit time : 
package model

import "gotest/model/mdata"

func (s *TestA) GetId() uint64 {
	return s.id
}
func (s *TestA) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(0)
}
func (s *TestA) GetA() int64 {
	return s.a
}
func (s *TestA) SetA(v int64) {
	s.a = v
	s.UpdateDirty(1)
}
func (s *TestA) GetB() int32 {
	return s.b
}
func (s *TestA) SetB(v int32) {
	s.b = v
	s.UpdateDirty(2)
}
func (s *TestA) CleanDirty() {
	s.DirtyModel.CleanDirty()
}
func (s *TestB) GetId() uint64 {
	return s.id
}
func (s *TestB) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(0)
}
func (s *TestB) GetM() string {
	return s.m
}
func (s *TestB) SetM(v string) {
	s.m = v
	s.UpdateDirty(1)
}
func (s *TestB) GetN() *TestA {
	return s.n
}
func (s *TestB) SetN(v *TestA) {
	s.n = v
	s.UpdateDirty(2)
	if v != nil {
		v.SetSelfDirtyIdx(2, s.UpdateDirty)
	}
}
func (s *TestB) GetC() *mdata.MList[*TestA] {
	return s.c
}
func (s *TestB) SetC(v *mdata.MList[*TestA]) {
	s.c = v
	s.UpdateDirty(3)
	if v != nil {
		v.SetSelfDirtyIdx(3, s.UpdateDirty)
	}
}
func (s *TestB) GetD() *mdata.MMap[uint64, *TestA] {
	return s.d
}
func (s *TestB) SetD(v *mdata.MMap[uint64, *TestA]) {
	s.d = v
	s.UpdateDirty(4)
	if v != nil {
		v.SetSelfDirtyIdx(4, s.UpdateDirty)
	}
}
func (s *TestB) CleanDirty() {
	s.DirtyModel.CleanDirty()
	if s.n != nil {
		s.n.CleanDirty()
	}
	if s.c != nil {
		s.c.CleanDirty()
	}
	if s.d != nil {
		s.d.CleanDirty()
	}
}
func (s *TestC) GetId() uint64 {
	return s.id
}
func (s *TestC) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(0)
}
func (s *TestC) GetX() string {
	return s.x
}
func (s *TestC) SetX(v string) {
	s.x = v
	s.UpdateDirty(1)
}
func (s *TestC) GetY() *TestB {
	return s.y
}
func (s *TestC) SetY(v *TestB) {
	s.y = v
	s.UpdateDirty(2)
	if v != nil {
		v.SetSelfDirtyIdx(2, s.UpdateDirty)
	}
}
func (s *TestC) CleanDirty() {
	s.DirtyModel.CleanDirty()
	if s.y != nil {
		s.y.CleanDirty()
	}
}
