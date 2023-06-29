package utils

type Snapshot struct {
	DeleteSet   DeleteSet
	StateVector map[int]int
}

func NewSnapshot(ds DeleteSet, stateMap map[int]int) {

}

func (s Snapshot) Equal(other *Snapshot) bool {
	if other == nil {
		return false
	}
	var ds1 = s.DeleteSet.Clients
	var ds2 = other.DeleteSet.Clients
	var sv1 = s.StateVector
	var sv2 = other.StateVector
	// todo 取指针对比可能有问题，可能需要修复
	if &ds1 != &ds2 || &sv1 != &sv2 {
		return false
	}
	for k, v := range sv1 {
		v2, ok := sv2[k]
		if !ok {
			return false
		}
		if v2 != v {
			return false
		}
	}
	for k, v := range ds1 {
		i, ok := ds2[k]
		if !ok {
			return false
		}
		if len(i) != len(v) {
			return false
		}
		for i2, item := range i {
			var dsItem2 = v[i2]
			if item.Clock != dsItem2.Clock || item.Length != dsItem2.Length {
				return false
			}
		}
	}
	return true
}

func (s Snapshot) EncodeSnapshotV2() []byte {
	return nil
}
