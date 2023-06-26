package model

import "YJS-GO/types"

// parent:Y.AbstractType|null
// insert(index:number, content:Array<object|boolean|Array|string|number|null|Uint8Array|Y.Type>)
// 在索引 处插入内容。请注意，内容是一个元素数组。即array.insert(0, [1])拼接list，在0位置插入1。
// push(Array<Object|boolean|Array|string|number|null|Uint8Array|Y.Type>)
// unshift(Array<Object|boolean|Array|string|number|null|Uint8Array|Y.Type>)
// delete(index:number, length:number)
// get(index:number)
// slice(start:number, end:number):Array<Object|boolean|Array|string|number|null|Uint8Array|Y.Type>
// 检索一系列内容
// length:number
//
// forEach(function(value:object|boolean|Array|string|number|null|Uint8Array|Y.Type,
// index:number, array: Y.Array))
//
// map(function(T, number, YArray):M):Array<M>
// toArray():Array<object|boolean|Array|string|number|null|Uint8Array|Y.Type>
// 将此 YArray 的内容复制到新数组。
// toJSON():Array<Object|boolean|Array|string|number|null>
// 将此 YArray 的内容复制到新数组。它使用它们的方法将所有子类型转换为 JSON toJSON。
// [Symbol.Iterator]
// 返回一个 YArray Iterator，它包含数组中每个索引的值。
// for (let value of yarray) { .. }
// observe(function(YArrayEvent, Transaction):void)
// 向该类型添加一个事件监听器，每次修改该类型时都会同步调用该事件监听器。如果在事件监听器中修改了该类型，当前事件监听器返回后会再次调用该事件监听器。
// unobserve(function(YArrayEvent, Transaction):void)
// observe从此类型中 删除事件侦听器。
// observeDeep(function(Array<YEvent>, Transaction):void)
// 向此类型添加一个事件侦听器，每次修改此类型或其任何子级时都会同步调用该事件侦听器。如果在事件监听器中修改了该类型，当前事件监听器返回后会再次调用该事件监听器。事件侦听器接收由其自身或其任何子级创建的所有事件。
// unobserveDeep(function(Array<YEvent>, Transaction):void)
// observeDeep从此类型中 删除事件侦听器。

type YArray struct {
	parent types.AbstractType
}

type YArrayEvent struct {
}
type Transaction struct {
}

func (*YArray) insert(index int, content YArray) {

}

// 在索引 处插入内容。请注意，内容是一个元素数组。即array.insert(0, [1])拼接list，在0位置插入1。
func (*YArray) push() {

}
func (*YArray) unshift() {

}
func (*YArray) delete(index int, length int) {

}
func (*YArray) get(index int) {

}
func (*YArray) slice(start int, end int) *YArray {
	return nil
}

// func (*YArray) // 检索一系列内容{
// }
func (*YArray) length() int {
	return 0
}

// forEach(function(value:object|boolean|Array|string|number|null|Uint8Array|Y.Type,
// index int, array: Y.Array))
func (*YArray) forEach(f func(abstractType types.AbstractType), index int) int {
	return 0
}

// map(function(T, number, YArray):M):Array<M>
func (*YArray) toArray() *YArray {
	return nil
}

// 将此 YArray 的内容复制到新数组。
func (*YArray) toJSON() *YArray {
	return nil
}

// 将此 YArray 的内容复制到新数组。它使用它们的方法将所有子类型转换为 JSON toJSON。
// [Symbol.Iterator]
// 返回一个 YArray Iterator，它包含数组中每个索引的值。
// for (let value of yarray) { .. }
func (*YArray) observe(a YArrayEvent, t Transaction) {

}

// 向该类型添加一个事件监听器，每次修改该类型时都会同步调用该事件监听器。如果在事件监听器中修改了该类型，当前事件监听器返回后会再次调用该事件监听器。
func (*YArray) unobserve(a YArrayEvent, t Transaction) {}

// observe从此类型中 删除事件侦听器。
func (*YArray) observeDeep(a []YEvent, t Transaction) {}

// 向此类型添加一个事件侦听器，每次修改此类型或其任何子级时都会同步调用该事件侦听器。如果在事件监听器中修改了该类型，当前事件监听器返回后会再次调用该事件监听器。事件侦听器接收由其自身或其任何子级创建的所有事件。
func (*YArray) unobserveDeep(a []YEvent, t Transaction) {}

// observeDeep从此类型中 删除事件侦听器。
