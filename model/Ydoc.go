package model

type YDoc struct {
	text  Text
	array Array[string, AbstractType]
}

type Text string

type Array[K string, V any] map[K]V

// 对共享文档应用文档更新。您可以选择指定 transactionOrigin将存储在 transaction.origin 和上ydoc.on('update', (update, origin) => ..)。
func (d *YDoc) applyUpdate() {

}

func (d *YDoc) encodeStateAsUpdate() {

}

func (d *YDoc) encodeStateVector() Uint8Array {
	return nil
}

func (d *YDoc) mergeUpdates([]Uint8Array) {

}

func (d *YDoc) encodeStateVectorFromUpdate(array Uint8Array) Uint8Array {
	return nil
}

// 将缺失的差异编码为另一个更新消息。此功能的工作方式类似于Y.encodeStateAsUpdate(ydoc, stateVector)但适用于更新。
func (d *YDoc) diffUpdate(update Uint8Array, stateVector Uint8Array) Uint8Array {
	return nil
}

func (d *YDoc) convertUpdateFormatV1ToV2() {

}

func (d *YDoc) convertUpdateFormatV2ToV1() {

}
