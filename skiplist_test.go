// Copyright (c) 2012, Glenn Brown.  All rights reserved.  See LICENSE.

package skiplist

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////

func TestT(t *testing.T) {
	t.Parallel()
	s := skiplist(1, 20)
	i := 1
	for e := s.Front(); e != nil; e = e.Next() {
		if e.Key().(int) != i || e.Value.(int) != 2*i {
			t.Fail()
		}
		i++
	}
}

func TestElement_Key(t *testing.T) {
	t.Parallel()
	e := skiplist(1, 3).Front()
	for i := 1; i <= 3; i++ {
		if e == nil || e.Key().(int) != i {
			t.Fail()
		}
		e = e.Next()
	}
}

func TestElement_Next(t *testing.T) {
	t.Parallel()
	e := skiplist(1, 3).Front()
	if e.Next().Key().(int) != 2 || e.Next().Next().Key().(int) != 3 {
		t.Fail()
	}
}

func TestElement_String(t *testing.T) {
	t.Parallel()
	if fmt.Sprint(skiplist(1, 2).Front()) != "1:2" {
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	// Verify that independent random number generators are used.
	s, s1 := New(), New()
	for i := 0; i < 32; i++ {
		s.Insert(i, i)
		s1.Insert(i, i)
	}
	v := s.visualization()
	v1 := s1.visualization()
	if v != v1 {
		t.Error("Not reproducible.")
	}
}

func TestNewDescending(t *testing.T) {
	t.Parallel()
	l := NewDescending().Insert(1, 1).Insert(2, 2).Insert(3, 3)
	if l.ElementN(0).Value.(int) != 3 || l.ElementN(1).Value.(int) != 2 || l.ElementN(2).Value.(int) != 1 {
		t.Fail()
	}
}

func TestT_Front(t *testing.T) {
	t.Parallel()
	s := skiplist(1, 3)
	if s.Front().Key().(int) != 1 {
		t.Fail()
	}
}

func TestT_Insert(t *testing.T) {
	t.Parallel()
	if skiplist(1, 10).String() != "{1:2 2:4 3:6 4:8 5:10 6:12 7:14 8:16 9:18 10:20}" {
		t.Fail()
	}
}

func TestT_Get(t *testing.T) {
	t.Parallel()
	l := skiplist(0, 7)
	if l.Get(0).(int) != 0 || l.Get(4).(int) != 8 || l.Get(7).(int) != 14 {
		t.Fail()
	}
}

func TestT_GetOk(t *testing.T) {
	t.Parallel()
	l := skiplist(1, 3)
	v, ok := l.GetOk(0)
	if nil != v || false != ok {
		t.Fail()
	}
	v, ok = l.GetOk(1)
	if 2 != v.(int) || true != ok {
		t.Fail()
	}
	v, ok = l.GetOk(3)
	if 6 != v.(int) || true != ok {
		t.Fail()
	}
	v, ok = l.GetOk(4)
	if nil != v || false != ok {
		t.Fail()
	}
}

func TestT_GetAll(t *testing.T) {
	t.Parallel()
	l := skiplist(1, 3).Insert(2, 3).Insert(2, 5)
	a := l.GetAll(0)
	if 0 != len(a) {
		t.Fail()
	}
	a = l.GetAll(1)
	if 1 != len(a) || 2 != a[0] {
		t.Fail()
	}
	a = l.GetAll(2)
	if 3 != len(a) || 5 != a[0] || 3 != a[1] || 4 != a[2] {
		t.Fail()
	}
	a = l.GetAll(3)
	if 1 != len(a) || 6 != a[0] {
		t.Fail()
	}
	a = l.GetAll(4)
	if 0 != len(a) {
		t.Fail()
	}
}

func TestT_Set(t *testing.T) {
	t.Parallel()
	l := skiplist(1, 3)
	l.Set(2, 2)
	a := l.GetAll(2)
	if 1 != len(a) || a[0].(int) != 2 {
		t.Fail()
	}
}

func TestT_Remove(t *testing.T) {
	t.Parallel()
	s := skiplist(0, 10)
	if s.Remove(-1) != nil || s.Remove(11) != nil {
		t.Error("Removing nonexistant key should fail.")
	}
	for i := range rand.Perm(11) {
		e := s.Remove(i)
		if e == nil {
			t.Error("nil")
		}
		if e.Key().(int) != i {
			t.Error("bad key")
		}
		if e.Value.(int) != 2*i {
			t.Error("bad value")
		}
	}
	if s.Len() != 0 {
		t.Error("nonzero len")
	}
}

func TestT_RemoveElement(t *testing.T) {
	t.Parallel()
	l := skiplist(0, 10)
	for i := 0; i <= 10; i += 2 {
		l.RemoveElement(l.Element(i))
	}
	if fmt.Sprintf(l.String()) != "{1:2 3:6 5:10 7:14 9:18}" {
		t.Fail()
	}
}

func TestT_RemoveN(t *testing.T) {
	t.Parallel()
	s := skiplist(0, 10)
	keys := rand.Perm(11)
	cnt := 11
	for _, key := range keys {
		found, pos := s.ElementPos(key)
		t.Logf("Removing key=%v at pos=%v", key, pos)
		t.Log(key, found, pos)
		t.Log("\n" + s.visualization())
		e := s.RemoveN(pos)
		if e == nil {
			t.Error("nil returned")
		} else if found != e {
			t.Error("Wrong removed")
		} else if e.Key().(int) != key {
			t.Error("bad Key()")
		} else if e.Value.(int) != 2*key {
			t.Error("bad Value")
		}
		cnt--
		l := s.Len()
		if l != cnt {
			t.Error("bad Len()=", l, "!=", cnt)
		}
	}
}

func TestT_ElementPos(t *testing.T) {
	t.Parallel()
	l := skiplist(1, 10)
	for i := 1; i < 10; i++ {
		e, pos := l.ElementPos(i)
		if e == nil || pos != i-1 {
			t.Fail()
		}
	}
	e, pos := l.ElementPos(0)
	if e != nil || pos != -1 {
		t.Fail()
	}
	e, pos = l.ElementPos(11)
	if e != nil || pos != -1 {
		t.Fail()
	}
}

func TestT_Element_forward(t *testing.T) {
	t.Parallel()
	s := skiplist(0, 9)
	for i := s.Len() - 1; i >= 0; i-- {
		e, pos := s.ElementPos(i)
		if e == nil {
			t.Error("nil")
		} else if e != s.ElementN(pos) {
			t.Error("bad pos")
		} else if e.Key().(int) != i {
			t.Error("bad Key")
		} else if e.Value.(int) != 2*i {
			t.Error("bad Value")
		}
	}
}

func TestT_Pos(t *testing.T) {
	t.Parallel()
	l := skiplist(0, 4)
	for i := 0; i < 5; i++ {
		if l.Pos(i) != i {
			t.Fail()
		}
	}
	if l.Pos(-2) != -1 || l.Pos(5) != -1 {
		t.Fail()
	}
}

func TestT_Len(t *testing.T) {
	t.Parallel()
	s := skiplist(0, 4)
	if s.Len() != 5 {
		t.Fail()
	}
}

func TestT_ElementN(t *testing.T) {
	t.Parallel()
	s := skiplist(0, 9)
	for i := s.Len() - 1; i >= 0; i-- {
		e := s.ElementN(i)
		if e == nil {
			t.Error("nil")
		} else if e.Key().(int) != i {
			t.Error("bad Key")
		} else if e.Value.(int) != 2*i {
			t.Error("bad Value")
		}
	}
}

func TestBuiltins(t *testing.T) {
	t.Parallel()

	// Create high and low variables for each ordered builtin type.

	f32a, f32b := float32(1.0), float32(0.0)
	f64a, f64b := float64(1.0), float64(0.0)
	i16a, i16b := int16(1), int16(0)
	i32a, i32b := int32(1), int32(0)
	i64a, i64b := int64(1), int64(0)
	i8_a, i8_b := int8(1), int8(0)
	i__a, i__b := int(1), int(0)
	sl_a, sl_b := []byte{1}, []byte{0}
	stra, strb := "1", "0"
	u16a, u16b := uint16(1), uint16(0)
	u32a, u32b := uint32(1), uint32(0)
	u64a, u64b := uint64(1), uint64(0)
	u8_a, u8_b := uint8(1), uint8(0)
	u__a, u__b := uint(1), uint(0)
	up_a, up_b := uintptr(1), uintptr(0)

	// Insert pairs in a map and verify the large is in position 1.

	if New().Set(f32a, 1).Set(f32b, 2).Pos(f32a) != 1 {
		t.Error("float32")
	}
	if New().Set(f64a, 1).Set(f64b, 2).Pos(f64a) != 1 {
		t.Error("float64")
	}
	if New().Set(i16a, 1).Set(i16b, 2).Pos(i16a) != 1 {
		t.Error("int16")
	}
	if New().Set(i32a, 1).Set(i32b, 2).Pos(i32a) != 1 {
		t.Error("int32")
	}
	if New().Set(i64a, 1).Set(i64b, 2).Pos(i64a) != 1 {
		t.Error("int64")
	}
	if New().Set(i8_a, 1).Set(i8_b, 2).Pos(i8_a) != 1 {
		t.Error("int8")
	}
	if New().Set(i__a, 1).Set(i__b, 2).Pos(i__a) != 1 {
		t.Error("int")
	}
	if New().Set(sl_a, 1).Set(sl_b, 2).Pos(sl_a) != 1 {
		t.Error("[]byte")
	}
	if New().Set(stra, 1).Set(strb, 2).Pos(stra) != 1 {
		t.Error("string")
	}
	if New().Set(u16a, 1).Set(u16b, 2).Pos(u16a) != 1 {
		t.Error("uint16")
	}
	if New().Set(u32a, 1).Set(u32b, 2).Pos(u32a) != 1 {
		t.Error("uint32")
	}
	if New().Set(u64a, 1).Set(u64b, 2).Pos(u64a) != 1 {
		t.Error("uint64")
	}
	if New().Set(u8_a, 1).Set(u8_b, 2).Pos(u8_a) != 1 {
		t.Error("uint8")
	}
	if New().Set(u__a, 1).Set(u__b, 2).Pos(u__a) != 1 {
		t.Error("uint")
	}
	if New().Set(up_a, 1).Set(up_b, 2).Pos(up_a) != 1 {
		t.Error("uintptr")
	}

	// Insert pairs in a map and verify the large is in position 1.

	if NewDescending().Set(f32a, 1).Set(f32b, 2).Pos(f32b) != 1 {
		t.Error("float32")
	}
	if NewDescending().Set(f64a, 1).Set(f64b, 2).Pos(f64b) != 1 {
		t.Error("float64")
	}
	if NewDescending().Set(i16a, 1).Set(i16b, 2).Pos(i16b) != 1 {
		t.Error("int16")
	}
	if NewDescending().Set(i32a, 1).Set(i32b, 2).Pos(i32b) != 1 {
		t.Error("int32")
	}
	if NewDescending().Set(i64a, 1).Set(i64b, 2).Pos(i64b) != 1 {
		t.Error("int64")
	}
	if NewDescending().Set(i8_a, 1).Set(i8_b, 2).Pos(i8_b) != 1 {
		t.Error("int8")
	}
	if NewDescending().Set(i__a, 1).Set(i__b, 2).Pos(i__b) != 1 {
		t.Error("int")
	}
	if NewDescending().Set(sl_a, 1).Set(sl_b, 2).Pos(sl_b) != 1 {
		t.Error("[]byte")
	}
	if NewDescending().Set(stra, 1).Set(strb, 2).Pos(strb) != 1 {
		t.Error("string")
	}
	if NewDescending().Set(u16a, 1).Set(u16b, 2).Pos(u16b) != 1 {
		t.Error("uint16")
	}
	if NewDescending().Set(u32a, 1).Set(u32b, 2).Pos(u32b) != 1 {
		t.Error("uint32")
	}
	if NewDescending().Set(u64a, 1).Set(u64b, 2).Pos(u64b) != 1 {
		t.Error("uint64")
	}
	if NewDescending().Set(u8_a, 1).Set(u8_b, 2).Pos(u8_b) != 1 {
		t.Error("uint8")
	}
	if NewDescending().Set(u__a, 1).Set(u__b, 2).Pos(u__b) != 1 {
		t.Error("uint")
	}
	if NewDescending().Set(up_a, 1).Set(up_b, 2).Pos(up_b) != 1 {
		t.Error("uintptr")
	}
}

////////////////////////////////////////////////////////////////
// Examples
////////////////////////////////////////////////////////////////

func Example() {
	// Create a skiplist and add some entries:
	s := New().Set("one", "un").Set("two", nil).Set("three", "trois")

	// Retrieve a mapping:
	fmt.Println(1, s.Get("two"))

	// Replace a mapping:
	s.Set("two", "deux")

	// Print the skiplist:
	fmt.Println(2, s)

	// Add more than one value for a key, even of different value-type:
	s.Insert("three", 3)

	// Retrieve all values for the key:
	fmt.Println(3, s.GetAll("three"))

	// Or just the youngest:
	fmt.Println(4, s.Get("three"))

	// Iterate over all values in the map:
	fmt.Print(5)
	for e := s.Front(); nil != e; e = e.Next() {
		fmt.Print(" ", e.Key(), "->", e.Value)
	}
	fmt.Println()

	// Pop the first entry:
	s.RemoveN(0)

	// Pop the last entry:
	s.RemoveN(s.Len() - 1)
	fmt.Println(6, s)

	// Output:
	// 1 <nil>
	// 2 {one:un three:trois two:deux}
	// 3 [3 trois]
	// 4 3
	// 5 one->un three->3 three->trois two->deux
	// 6 {three:3 three:trois}
}

// This example demonstrates iteration over all list elements.
func ExampleElement_Next() {
	s := New().Insert(0, 0).Insert(1, 1).Insert(1, 2).Insert(2, 4)

	// Efficiently iterate over all entries:
	fmt.Print("A")
	for e := s.Front(); e != nil; e = e.Next() {
		fmt.Print(" ", e)
	}
	fmt.Println()

	// Efficiently iterate over entries for a single key:
	fmt.Print("B")
	for e := s.Element(1); e != nil && e.Key().(int) == 1; e = e.Next() {
		fmt.Print(" ", e)
	}

	// Output:
	// A 0:0 1:2 1:1 2:4
	// B 1:2 1:1
}

func ExampleElement_String() {
	fmt.Println(New().Set("key1", "value1").ElementN(0))
	// Output: key1:value1
}

func ExampleT_GetAll() {
	s := New().Insert(0, 0).Insert(1, 1).Insert(1, 2).Insert(2, 4)

	// Conveniently iterate over values for a single key:
	for _, ee := range s.GetAll(1) {
		fmt.Print(" ", ee)
	}
	// Output: 2 1
}

func ExampleT_String() {
	skip := New().Insert(1, 10).Insert(2, 20).Insert(3, 30)
	fmt.Println(skip)
	// Output: {1:10 2:20 3:30}
}

// One may Remove() during iteration.
func ExampleT_Remove() {
	skip := New().Insert(1, 10).Insert(2, 20).Insert(3, 30)
	for e := skip.Front(); nil != e; e = e.Next() {
		fmt.Println(skip.Remove(e.Key()))
	}
	// Output:
	// 1:10
	// 2:20
	// 3:30
}

// One may RemoveElement() during iteration.
func ExampleT_RemoveElement() {
	skip := New().Insert(1, 10).Insert(2, 20).Insert(3, 30)
	for e := skip.Front(); nil != e; e = e.Next() {
		fmt.Println(skip.RemoveElement(e))
	}
	// Output:
	// 1:10
	// 2:20
	// 3:30
}

func TestVisualization(t *testing.T) {
	s := New()
	for i := 0; i < 23; i++ {
		s.Insert(i, i)
	}
	v := s.visualization()
	expected := "" +
		"L4 |---------------------------------------------------------------------->/\n" +
		"L3 |---------------------->|---->|---------------------------------------->/\n" +
		"L2 |---------------------->|->|->|---------------->|---------------------->/\n" +
		"L1 |------->|------------->|->|->|->|---->|---->|->|---->|---------->|---->/\n" +
		"L0 |->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->/\n" +
		"      0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  1  1  1  1  1  1  1\n" +
		"      0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f  0  1  2  3  4  5  6"
	if v != expected {
		t.Error(v, "\n!=\n", expected)
	}
}

////////////////////////////////////////////////////////////////
// Benchmarks
////////////////////////////////////////////////////////////////

func BenchmarkT_Insert_forward(b *testing.B) {
	b.StopTimer()
	s := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Insert(i, i)
	}
}

func BenchmarkT_Insert_reverse(b *testing.B) {
	b.StopTimer()
	s := New()
	b.StartTimer()
	for i := b.N - 1; i >= 0; i-- {
		s.Insert(i, i)
	}
}

func BenchmarkT_Insert_shuffle(b *testing.B) {
	b.StopTimer()
	a := rand.Perm(b.N)
	s := New()
	b.StartTimer()
	for i, key := range a {
		s.Insert(key, i)
	}
}

func BenchmarkT_Element_forward(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := b.N - 1; i >= 0; i-- {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Element(i)
	}
}

func BenchmarkT_Element_reverse(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := 0; i < b.N; i++ {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := b.N - 1; i >= 0; i-- {
		s.Element(i)
	}
}

func BenchmarkT_Element_shuffle(b *testing.B) {
	b.StopTimer()
	a := rand.Perm(b.N)
	s := skiplist(0, b.N-1)
	b.StartTimer()
	for _, key := range a {
		s.Element(key)
	}
}

func BenchmarkT_ElementN_forward(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := b.N - 1; i >= 0; i-- {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.ElementN(i)
	}
}

func BenchmarkT_ElementN_reverse(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := 0; i < b.N; i++ {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := b.N - 1; i >= 0; i-- {
		s.ElementN(i)
	}
}

func BenchmarkT_ElementN_shuffle(b *testing.B) {
	b.StopTimer()
	a := rand.Perm(b.N)
	s := skiplist(0, b.N-1)
	b.StartTimer()
	for _, key := range a {
		s.ElementN(key)
	}
}

func BenchmarkT_Remove_forward(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := b.N - 1; i >= 0; i-- {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Remove(i)
	}
}

func BenchmarkT_Remove_reverse(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := 0; i < b.N; i++ {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := b.N - 1; i >= 0; i-- {
		s.Remove(i)
	}
}

func BenchmarkT_Remove_shuffle(b *testing.B) {
	b.StopTimer()
	a := rand.Perm(b.N)
	s := skiplist(0, b.N-1)
	b.StartTimer()
	for _, key := range a {
		s.Remove(key)
	}
}

func BenchmarkT_RemoveN_head(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := b.N - 1; i >= 0; i-- {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.RemoveN(0)
	}
}

func BenchmarkT_RemoveN_tail(b *testing.B) {
	b.StopTimer()
	s := New()
	for i := 0; i < b.N; i++ {
		s.Insert(i, i)
	}
	b.StartTimer()
	for i := b.N - 1; i >= 0; i-- {
		s.RemoveN(i)
	}
}

func BenchmarkT_RemoveN_mid(b *testing.B) {
	b.StopTimer()
	s := skiplist(0, b.N-1)
	b.StartTimer()
	for i := b.N - 1; i >= 0; i-- {
		s.RemoveN(i / 2)
	}
}

////////////////////////////////////////////////////////////////
// Utility functions
////////////////////////////////////////////////////////////////

// Create a skiplist with each key in [min,max].
//
func skiplist(min, max int) *T {
	s := New()
	for _, v := range rand.Perm(max - min + 1) {
		s.Insert((min + v), 2*(min+v))
	}
	return s
}

// Create an arrow string like "|-->" that is cnt runes long.
//
func arrow(cnt int) (s string) {
	cnt *= 3
	switch {
	case cnt > 1:
		return "|" + strings.Repeat("-", cnt-2) + ">"
	case cnt == 1:
		return ">"
	}
	return "X"
}

// Create a visualization string like this:
//   Output:
//   L4 |---------------------------------------------------------------------->/
//   L3 |------------------------------------------->|------------------------->/
//   L2 |---------->|---------->|---------->|------->|---------------->|---->|->/
//   L1 |---------->|---------->|---------->|->|---->|->|->|->|------->|->|->|->/
//   L0 |->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->|->/
//         0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  1  1  1  1  1  1  1
//         0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f  0  1  2  3  4  5  6
//
func (l *T) visualization() (s string) {
	for level := len(l.links) - 1; level >= 0; level-- {
		s += fmt.Sprintf("L%d ", level)
		w := l.links[level].width
		s += arrow(w)
		for n := l.links[level].to; n != nil; n = n.links[level].to {
			w = n.links[level].width
			s += arrow(w)
		}
		s += "/\n"
	}
	s += "    "
	for n := l.links[0].to; n != nil; n = n.links[0].to {
		s += fmt.Sprintf("  %x", n.key.(int)>>4&0xf)
	}
	s += "\n    "
	for n := l.links[0].to; n != nil; n = n.links[0].to {
		s += fmt.Sprintf("  %x", n.key.(int)&0xf)
	}
	return string(s)
}
