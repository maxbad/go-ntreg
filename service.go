package ntreg

/*
#if defined(USE_NTREG)
#include <stdlib.h>
#include <stdint.h>
#include "ntreg.h"
#endif
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// CreateService 创建服务
func CreateService(serviceName string, registryFilePath string) error {
	// 转换go的字符串为C字符
	_serviceName := C.CString(serviceName)
	defer C.free(unsafe.Pointer(_serviceName))
	_registryFilePath := C.CString(registryFilePath)
	defer C.free(unsafe.Pointer(_registryFilePath))

	// 打开注册表文件
	_regHive := C.openHive(_registryFilePath, C.HMODE_RW)
	if _regHive == nil {
		return errors.New("error opening registry file")
	}
	defer C.closeHive(_regHive)

	// 判断注册表类型
	if C.int(_regHive._type) != C.HTYPE_SYSTEM {
		return errors.New("not a registry file or type error")
	}

	// 打开服务路径
	_servicePath := C.CString("ControlSet001\\Services")
	defer C.free(unsafe.Pointer(_servicePath))
	nkPos := C.trav_path(_regHive, C.int(0), _servicePath, C.REG_NONE)
	if nkPos == C.int(0) {
		return errors.New("failed to open service path")
	}

	// 创建服务名称
	if C.add_key(_regHive, nkPos+4, _serviceName) == nil {
		return errors.New("failed to create service name. The service may already exist")
	}

	// 打开服务名称路径
	nkPos = C.trav_path(_regHive, nkPos+4, _serviceName, C.REG_NONE)
	if nkPos == 0 {
		return errors.New("failed to open service name path")
	}

	// 写入Type值
	if errorCode := writeDwordValue(_regHive, nkPos+4, "Type", C.REG_DWORD, C.TPF_VK_ABS|C.TPF_EXACT, 1); errorCode != 0 {
		return errors.New(fmt.Sprintf("failed to write Type value #%d", errorCode))
	}

	// 写入Start值
	if errorCode := writeDwordValue(_regHive, nkPos+4, "Start", C.REG_DWORD, C.TPF_VK_ABS|C.TPF_EXACT, 0); errorCode != 0 {
		return errors.New(fmt.Sprintf("failed to write Start value #%d", errorCode))
	}

	// 写入ErrorControl值
	if errorCode := writeDwordValue(_regHive, nkPos+4, "ErrorControl", C.REG_DWORD, C.TPF_VK_ABS|C.TPF_EXACT, 1); errorCode != 0 {
		return errors.New(fmt.Sprintf("failed to write ErrorControl value #%d", errorCode))
	}

	// 写入Group值
	if errorCode := writeCharValue(_regHive, nkPos+4, "Group", C.REG_SZ, C.TPF_VK_ABS|C.TPF_EXACT, "NDIS Wrapper"); errorCode != 0 {
		return errors.New(fmt.Sprintf("failed to write Group value #%d", errorCode))
	}

	// 写入ImagePath值
	imagePath := fmt.Sprintf("System32\\\\drivers\\\\%s.sys", serviceName)
	if errorCode := writeCharValue(_regHive, nkPos+4, "ImagePath", C.REG_EXPAND_SZ, C.TPF_VK_ABS|C.TPF_EXACT, imagePath); errorCode != 0 {
		return errors.New(fmt.Sprintf("failed to write ImagePath value #%d", errorCode))
	}

	// 写出到磁盘
	if C.writeHive(_regHive) != 0 {
		return errors.New("failed to write to disk")
	}

	return nil
}

// DeleteService 删除服务
func DeleteService(serviceName string, registryFilePath string) error {
	// 转换go的字符串为C字符
	_serviceName := C.CString(serviceName)
	defer C.free(unsafe.Pointer(_serviceName))
	_registryFilePath := C.CString(registryFilePath)
	defer C.free(unsafe.Pointer(_registryFilePath))

	// 打开注册表文件
	_regHive := C.openHive(_registryFilePath, C.HMODE_RW)
	if _regHive == nil {
		return errors.New("error opening registry file")
	}
	defer C.closeHive(_regHive)

	// 判断注册表类型
	if C.int(_regHive._type) != C.HTYPE_SYSTEM {
		return errors.New("not a registry file or type error")
	}

	// 打开服务路径
	_servicePath := C.CString("ControlSet001\\Services")
	defer C.free(unsafe.Pointer(_servicePath))
	nkPos := C.trav_path(_regHive, C.int(0), _servicePath, C.REG_NONE)
	if nkPos == 0 {
		return errors.New("failed to open service path")
	}

	// 打开服务名称路径
	nkPosSub := C.trav_path(_regHive, nkPos+4, _serviceName, C.REG_NONE)
	if nkPosSub == 0 {
		return errors.New("failed to open service name path")
	}

	// 删除名称下面的子项
	C.del_allvalues(_regHive, nkPosSub+4)

	// 删除服务项
	if C.del_key(_regHive, nkPos+4, _serviceName) != 0 {
		return errors.New("failed to delete service item")
	}

	// 写出到磁盘
	if C.writeHive(_regHive) != 0 {
		return errors.New("failed to write to disk")
	}

	return nil
}

// 写入DWORD值的函数
func writeDwordValue(h *C.struct_hive, nkPos C.int, valueName string, valueType C.int, flags C.int, value uint32) int {
	_valueName := C.CString(valueName)
	defer C.free(unsafe.Pointer(_valueName))
	if C.add_value(h, nkPos, _valueName, valueType) == nil {
		return 1
	}
	if C.put_dword(h, nkPos, _valueName, flags, C.int(value)) == 0 {
		return 2
	}
	return 0
}

// 写入字符型值的函数
func writeCharValue(h *C.struct_hive, nkPos C.int, valueName string, valueType C.int, flags C.int, value string) int {
	_valueName := C.CString(valueName)
	defer C.free(unsafe.Pointer(_valueName))
	if C.add_value(h, nkPos, _valueName, valueType) == nil {
		return 1
	}
	_value := C.CString(value)
	defer C.free(unsafe.Pointer(_value))
	if C.put_char(h, nkPos, _valueName, valueType, flags, _value) == 0 {
		return 2
	}
	return 0
}
