// Code generated by "stringer -type itemType"; DO NOT EDIT.

package transform

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[itemError-0]
	_ = x[itemEOF-1]
	_ = x[itemNewLine-2]
	_ = x[itemSpecial-3]
	_ = x[itemIdent-4]
	_ = x[itemComment-5]
	_ = x[itemFileHeader-6]
	_ = x[itemTypeHeader-7]
	_ = x[itemString-8]
	_ = x[itemValue-9]
	_ = x[itemCommand-10]
	_ = x[itemWord-11]
	_ = x[itemKeyword-12]
}

const _itemType_name = "itemErroritemEOFitemNewLineitemSpecialitemIdentitemCommentitemFileHeaderitemTypeHeaderitemStringitemValueitemCommanditemWorditemKeyword"

var _itemType_index = [...]uint8{0, 9, 16, 27, 38, 47, 58, 72, 86, 96, 105, 116, 124, 135}

func (i itemType) String() string {
	if i < 0 || i >= itemType(len(_itemType_index)-1) {
		return "itemType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _itemType_name[_itemType_index[i]:_itemType_index[i+1]]
}