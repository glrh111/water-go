package hashtable

import (
	"unsafe"
	"fmt"
)

type HashKeyer interface {
	HashCode() int
	IsEqual(HashKeyer) bool
	fmt.Stringer
}

type IntHashKey struct {
	value int
	hashCode int
}

func NewIntHashKey(v int) *IntHashKey {
	return &IntHashKey{v, -1}
}

func (ihk *IntHashKey) HashCode() int { return ihk.value}

func (ihk *IntHashKey) IsEqual(ano HashKeyer) bool {
	if value, ok := ano.(*IntHashKey); ok {
		return value.value == ihk.value
	}
	return false
}

func (ihk *IntHashKey) String() string {
	return fmt.Sprintf("IntHashKey: [%v] HashCode: [%v]", ihk.value, ihk.HashCode())
}

func NewStringHashKey(s string) *StringHashKey {
	return &StringHashKey{s, -1}
}

type StringHashKey struct {
	value string
	hashCode int  // -1 by default
}

func (shk *StringHashKey) HashCode() int {
	var (
		x = hashSecret.prefix
		ls = len(shk.value)
		lp = ls

	)
	if shk.hashCode != -1 {
		return shk.hashCode
	}
	if ls == 0 {
		shk.hashCode = 0
		return 0
	}
	x ^= int(shk.value[0] << 7)
	lp -= 1
	for ; lp > 0; lp-- {
		x = (100000*x) ^ int(shk.value[ls-lp])
	}
	x ^= int(uintptr(unsafe.Pointer(shk)))
	x ^= hashSecret.suffix
	shk.hashCode = x
	return x
}

func (shk *StringHashKey) IsEqual(ano HashKeyer) bool {
	if value, ok := ano.(*StringHashKey); ok {
		return value.value == shk.value
	}
	return false
}

func (shk *StringHashKey) String() string {
	return fmt.Sprintf("StringHashKey: [%v] HashCode: [%v]", shk.value, shk.HashCode())
}


