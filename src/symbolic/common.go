package symbolic

//https://msdn.microsoft.com/en-us/library/cc232007.aspx
//https://gist.github.com/Perlmint/f9f0e37db163dd69317d
//https://github.com/golang/go/blob/master/src/os/os_windows_test.go

import (
	"syscall"
	//"encoding/binary"
	//"bytes"
	"unsafe"
	"strings"
	//"errors"
	"unicode/utf16"
)

func IsReparsePoint(path string) (ret bool, err error) {
	attrs, err := syscall.GetFileAttributes(syscall.StringToUTF16Ptr(path));
	if err != nil {
		return false, err
	}
	return attrs & syscall.FILE_ATTRIBUTE_DIRECTORY > 0, nil
}

/*
reparseType -> syscall.IO_REPARSE_TAG_SYMLINK 符号链接
	       IO_REPARSE_TAG_MOUNT_POINT     junction
 */
func getReparseTarget(handle syscall.Handle) (target string, reparseType int, err error) {

	var bytesReturned uint32;

	var outBuffer MountPointReparseBuffer;

	//junction point => IO_REPARSE_TAG_MOUNT_POINT (0xA0000003).
	//symbolink point => SYMBOLIC_LINK_FLAG_DIRECTORY  (0xA000000C)
	err = syscall.DeviceIoControl(handle, syscall.FSCTL_GET_REPARSE_POINT,
		nil, 0, (*byte)(unsafe.Pointer(&outBuffer)), uint32(unsafe.Sizeof(outBuffer)), &bytesReturned, nil);

	if err == nil {

		if outBuffer.ReparseTag == syscall.IO_REPARSE_TAG_SYMLINK {

			var symbolicLink *SymbolicLinkReparseBuffer;
			symbolicLink = (*SymbolicLinkReparseBuffer)(unsafe.Pointer(&outBuffer))

			target = symbolicLink.SubstituteName()
			reparseType = int(outBuffer.ReparseTag)

		} else if outBuffer.ReparseTag == IO_REPARSE_TAG_MOUNT_POINT {

			target = outBuffer.SubstituteName();
			reparseType = int(outBuffer.ReparseTag)
		}
		if (strings.HasPrefix(target, NonInterpretedPathPrefix)) {
			target = target[len(NonInterpretedPathPrefix):]
		}
	} else {
		target = ""
	}
	return
}

func openReparsePoint(reparsePoint string, accessMode uint32) (handle syscall.Handle, err error) {
	handle, err = syscall.CreateFile(
		syscall.StringToUTF16Ptr(reparsePoint),
		accessMode,
		syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS | syscall.FILE_FLAG_OPEN_REPARSE_POINT,
		0)
	return
}

func (r *SymbolicLinkReparseBuffer) PrintName() string {
	offset := r.PrintNameOffset / 2//微软官方文档中说要除了2
	length := r.PrintNameLength / 2
	return string(utf16.Decode(r.PathBuffer[offset:offset + length]))
}

func (r *SymbolicLinkReparseBuffer) SubstituteName() string {
	offset := r.SubstituteNameOffset / 2//微软官方文档中说要除了2
	length := r.SubstituteNameLength / 2
	return string(utf16.Decode(r.PathBuffer[offset:offset + length]))
}

func (r *MountPointReparseBuffer) PrintName() string {
	offset := r.PrintNameOffset / 2//微软官方文档中说要除了2
	length := r.PrintNameLength / 2
	return string(utf16.Decode(r.PathBuffer[offset:offset + length]))
}

func (r *MountPointReparseBuffer) SubstituteName() string {
	offset := r.SubstituteNameOffset / 2//微软官方文档中说要除了2
	length := r.SubstituteNameLength / 2
	return string(utf16.Decode(r.PathBuffer[offset:offset + length]))
}
