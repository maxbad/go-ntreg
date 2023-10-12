//go:build ntreg
// +build ntreg

package go_ntreg

/*
#cgo CFLAGS: -DUSE_NTREG
#if defined(USE_NTREG)
#include "ntreg.h"
#endif
*/
import "C"

// hive 定义，由 openHive() 分配，由 closeHive() 释放
// 包含状态数据，必须在所有函数中传递
type Hive struct {
	filename    string  // hive 文件名
	filedesc    int     // 文件描述符（仅当 state == OPEN 时有效）
	state       int     // hive 的当前状态
	htype       int     // hive 推荐的类型。注意：库加载时会自动猜测，但应用程序可以根据需要更改它
	pages       int     // 页面总数
	useblk      int     // 使用的块总数
	unuseblk    int     // 未使用的块总数
	usetot      int     // useblk 中的总字节数
	unusetot    int     // unuseblk 中的总字节数
	size        int     // hive 的大小（文件大小），包括 regf 头部
	rootofs     int     // 根节点的偏移量
	lastbin     int     // 最后一个 HBIN 偏移量
	endofs      int     // 第一个非 HBIN 页面的偏移量，从这里可以扩展
	nkindextype int16   // 根键使用的子键索引类型
	buffer      *C.char // 文件的原始内容
}

// Hive open modes
const (
	HMODE_RW        = 0
	HMODE_RO        = 0x1
	HMODE_OPEN      = 0x2
	HMODE_DIRTY     = 0x4
	HMODE_NOALLOC   = 0x8  // Don't allocate new blocks
	HMODE_NOEXPAND  = 0x10 // Don't expand file with new hbin
	HMODE_DIDEXPAND = 0x20 // File has been expanded
	HMODE_VERBOSE   = 0x1000
	HMODE_TRACE     = 0x2000
	HMODE_INFO      = 0x4000 // Show some info on open and close
)

func OpenHive(filename string, mode int) *Hive {
	_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(filename))
	_mode := C.int(mode)
	return C.openHive(_filename, _mode)
}
