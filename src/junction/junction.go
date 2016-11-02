package junction

//https://msdn.microsoft.com/en-us/library/cc232007.aspx
//https://gist.github.com/Perlmint/f9f0e37db163dd69317d
//https://github.com/golang/go/blob/master/src/os/os_windows_test.go

import (
	"directory"
	"os"
	"syscall"
	"path/filepath"
	//"encoding/binary"

	//"bytes"
	"unsafe"
	"strings"
	"errors"
)




/// Creates a junction point from the specified directory to the specified target directory.
/// Only works on NTFS.
func CreateJunction(junctionPoint string, targetDir string, overwrite bool) (result bool, err error) {
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return false, err
	}
	var exist bool
	if exist, err = directory.DirectoryExist(junctionPoint); err != nil {
		return false, err
	}

	if exist {
		if !overwrite {
			return false, errors.New("Directory already exists and overwrite parameter is false.")
		}
	} else {
		if err = os.MkdirAll(junctionPoint, 0777); err != nil {
			return false, err
		}
	}

	var handle syscall.Handle
	handle, err = openReparsePoint(junctionPoint, syscall.GENERIC_WRITE)

	defer syscall.CloseHandle(handle)

	substituteName := syscall.StringToUTF16(NonInterpretedPathPrefix + targetDir)

	printName := syscall.StringToUTF16(targetDir)

	reparseDataBuffer := MountPointReparseBuffer{}

	reparseDataBuffer.ReparseTag = IO_REPARSE_TAG_MOUNT_POINT


	//加上8个字节的 SubstituteNameOffset,SubstituteNameLength,PrintNameOffset,PrintNameLength
	//官方文档  Mount Point Reparse Data Buffer
	//This value is the length of the data starting at the SubstituteNameOffset field (or the size of the PathBuffer field, in bytes, plus 8).
	//https://msdn.microsoft.com/en-us/library/cc232007.aspx
	reparseDataBuffer.ReparseDataLength = uint16((len(substituteName) + len(printName)) * 2 + 8)

	reparseDataBuffer.Reserved = 0

	reparseDataBuffer.SubstituteNameOffset = 0

	//减1是去掉最后的 \0
	reparseDataBuffer.SubstituteNameLength = uint16((len(substituteName)) * 2 - 2)

	reparseDataBuffer.PrintNameOffset = uint16(len(substituteName) * 2)

	//减1是去掉最后的 \0
	reparseDataBuffer.PrintNameLength = uint16((len(printName)) * 2 - 2)

	var i int
	for i = 0; i < len(substituteName); i++ {
		reparseDataBuffer.PathBuffer[i] = substituteName[i]
	}
	var j, lenth int
	j = 0
	lenth = len(printName)
	for i = len(substituteName); j < lenth; j++ {
		reparseDataBuffer.PathBuffer[i + j] = printName[j]
	}

	var bytesReturned uint32;

	////第三个参数 inBufferSize = sizeof(REPARSE_DATA_BUFFER_HEADER)=8 的header
	//即 4byte的ReparseTag+2byte的ReparseDataLength+2byte的Reserved
	err = syscall.DeviceIoControl(handle, FSCTL_SET_REPARSE_POINT,
		(*byte)(unsafe.Pointer(&reparseDataBuffer)), uint32(reparseDataBuffer.ReparseDataLength + 8), nil, 0, &bytesReturned, nil);

	if err != nil {
		return false, err
	}

	return true, nil
}

/// Deletes a junction point at the specified source directory along with the directory itself.
/// Does nothing if the junction point does not exist.
/// Only works on NTFS.
func DeleteJunction(junctionPoint string) (result bool, err error) {
	var exist bool
	if exist, err = directory.DirectoryExist(junctionPoint); err != nil || exist == false {
		return false, errors.New("directory can not open or not exist.");
	}

	var handle syscall.Handle

	handle, err = openReparsePoint(junctionPoint, syscall.GENERIC_WRITE)

	defer syscall.CloseHandle(handle)

	if err != nil {
		return false, err;
	}

	reparseDataBuffer := MountPointReparseBuffer{}

	reparseDataBuffer.ReparseTag = IO_REPARSE_TAG_MOUNT_POINT
	reparseDataBuffer.ReparseDataLength = 0
	reparseDataBuffer.Reserved = 0

	var bytesReturned uint32;

	//第三个参数 inBufferSize = sizeof(REPARSE_DATA_BUFFER_HEADER)=8 的header
	//见 https://msdn.microsoft.com/en-us/library/ff552012.aspx
	//见 https://msdn.microsoft.com/en-us/library/cc232005.aspx
	//见https://msdn.microsoft.com/en-us/library/cc232006.aspx
	//见根目录下  NTFS Hard Links, Directory Junctions, and Windows Shortcuts.mhtml
	err = syscall.DeviceIoControl(handle, FSCTL_DELETE_REPARSE_POINT,
		(*byte)(unsafe.Pointer(&reparseDataBuffer)), 8, nil, 0, &bytesReturned, nil);

	if err != nil {
		return false, err;
	}

	err = os.RemoveAll(junctionPoint)
	if err != nil {
		return false, err;
	}
	return true, err
}

/// Determines whether the specified path exists and refers to a junction point.
func IsJunction(path string) bool {
	var exist bool
	var err error
	if exist, err = directory.DirectoryExist(path); err != nil || exist == false {
		return false
	}
	var handle syscall.Handle
	handle, err = openReparsePoint(path, syscall.GENERIC_READ)

	//安全关闭 句柄
	defer syscall.CloseHandle(handle)

	if err != nil {
		return false
	}

	_, err = internalJunctionGetTarget(handle)
	if err != nil {
		return false
	}

	return true
}

/// Gets the target of the specified junction point.
/// Only works on NTFS.
func GetJunctionTarget(junctionPoint string) (target string, err error) {

	var handle syscall.Handle
	handle, err = openReparsePoint(junctionPoint, syscall.GENERIC_READ)

	if err == nil {
		target, err = internalJunctionGetTarget(handle);
	}

	defer syscall.CloseHandle(handle)

	return

}

func internalJunctionGetTarget(handle syscall.Handle) (target string, err error) {

	var bytesReturned uint32;

	var outBuffer MountPointReparseBuffer;

	//junction point => IO_REPARSE_TAG_MOUNT_POINT (0xA0000003).
	//symbolink point => SYMBOLIC_LINK_FLAG_DIRECTORY  (0xA000000C)
	err = syscall.DeviceIoControl(handle, syscall.FSCTL_GET_REPARSE_POINT,
		nil, 0, (*byte)(unsafe.Pointer(&outBuffer)), uint32(unsafe.Sizeof(outBuffer)), &bytesReturned, nil);

	if err == nil {
		if outBuffer.ReparseTag != IO_REPARSE_TAG_MOUNT_POINT {
			target = ""
			err = errors.New("This Directory is not junction point.")
			return
		}

		target = outBuffer.SubstituteName();
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


