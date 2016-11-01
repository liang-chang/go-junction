package junction

//https://msdn.microsoft.com/en-us/library/cc232006.aspx
//https://gist.github.com/Perlmint/f9f0e37db163dd69317d
//https://github.com/golang/go/blob/master/src/os/os_windows_test.go

import (
	"directory"
	"os"
	"fmt"
	"syscall"
	"path/filepath"
	//"encoding/binary"

	//"bytes"
	"unsafe"
	"unicode/utf16"
	"strings"
	"errors"
)

func (r *MountPointReparseBuffer) PrintName() string {
	offset := r.PrintNameOffset / 2//微软官方文档中说要除了2
	length := r.PrintNameLength / 2
	return string(utf16.Decode(r.PathBuffer[offset:offset + length]))
}

func (r *MountPointReparseBuffer) SubstituteName() string {
	offset := r.SubstituteNameOffset / 2//微软官方文档中说要除了2
	length := r.SubstituteNameLength / 2
	//fmt.Println(string(utf16.Decode(r.PathBuffer[:])))
	return string(utf16.Decode(r.PathBuffer[offset:offset + length]))
}


/// Creates a junction point from the specified directory to the specified target directory.
/// Only works on NTFS.
func Create(junctionPoint string, targetDir string, overwrite bool) (result bool, err error) {
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

	fmt.Println("1=" + string(utf16.Decode(substituteName)))

	reparseDataBuffer := MountPointReparseBuffer{}

	reparseDataBuffer.ReparseTag = IO_REPARSE_TAG_MOUNT_POINT
	reparseDataBuffer.ReparseDataLength = 1048 * 2 + 12  //官方文档  the size of the PathBuffer field, in bytes, plus 12

	reparseDataBuffer.Reserved = 0

	reparseDataBuffer.SubstituteNameOffset = 0
	reparseDataBuffer.SubstituteNameLength = uint16(len(NonInterpretedPathPrefix + targetDir) * 2)

	reparseDataBuffer.PrintNameOffset = uint16((len(NonInterpretedPathPrefix + targetDir)+1) * 2)
	reparseDataBuffer.PrintNameLength = uint16(len(targetDir) * 2)
	//reparseDataBuffer.PathBuffer = targetDirBytes//targetDirBytes[0:len(targetDirBytes) - 1]//Array.Copy(targetDirBytes, reparseDataBuffer.PathBuffer, targetDirBytes.Length);
	//copy(reparseDataBuffer.PathBuffer[:], substituteName);

	for i := 0; i < len(substituteName); i++ {
		reparseDataBuffer.PathBuffer[i] = substituteName[i]
	}
	for i := 0; i < len(printName); i++ {
		reparseDataBuffer.PathBuffer[int(reparseDataBuffer.PrintNameOffset) / 2 + i] = printName[i]
	}

	fmt.Println("2=" + reparseDataBuffer.SubstituteName())
	fmt.Println("3=" + reparseDataBuffer.PrintName())

	fmt.Println(string(utf16.Decode(reparseDataBuffer.PathBuffer[:])))

	var bytesReturned uint32;

	err = syscall.DeviceIoControl(handle, FSCTL_SET_REPARSE_POINT,
		(*byte)(unsafe.Pointer(&reparseDataBuffer)), uint32(unsafe.Sizeof(reparseDataBuffer) ), nil, 0, &bytesReturned, nil);

	//+ unsafe.Offsetof(reparseDataBuffer.SubstituteNameOffset)

	if err != nil {
		return false, err
	}

	return true, nil
}

/// Deletes a junction point at the specified source directory along with the directory itself.
/// Does nothing if the junction point does not exist.
/// Only works on NTFS.
func Delete(junctionPoint string) (result bool, err error) {
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
	reparseDataBuffer.Reserved=0

	var bytesReturned uint32;

	//第三个参数 inBufferSize 是 8 _REPARSE_DATA_BUFFER 的header
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
func Exists(path string) bool {
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

	_, err = internalGetTarget(handle)
	if err != nil {
		return false
	}

	return true
}

/// Gets the target of the specified junction point.
/// Only works on NTFS.
func GetTarget(junctionPoint string) (target string, err error) {
	var handle syscall.Handle
	handle, err = openReparsePoint(junctionPoint, syscall.GENERIC_READ)

	if err == nil {
		target, err = internalGetTarget(handle);
	}

	defer syscall.CloseHandle(handle)

	return

}

func internalGetTarget(handle syscall.Handle) (target string, err error) {

	var bytesReturned uint32;

	var outBuffer MountPointReparseBuffer;

	err = syscall.DeviceIoControl(handle, syscall.FSCTL_GET_REPARSE_POINT,
		nil, 0, (*byte)(unsafe.Pointer(&outBuffer)), uint32(unsafe.Sizeof(outBuffer)), &bytesReturned, nil);

	if err == nil {
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


