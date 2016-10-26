package main

import (
	"directory"
	"config"
	//	"fmt"
	//	"regexp"
	//"log"
	//	"os"
	//	"os/exec"
	//	"path/filepath"
	//	"strings"
	//	"bytes"
	//	"strings"
	//	"bytes"
	//	"path/filepath"
	//	"os"
	//	"log"
	//"io/ioutil"
	//"fmt"

	"os"
	"fmt"
	"syscall"
	"unsafe"
	"path/filepath"
)

/// <summary>
/// The file or directory is not a reparse point.
/// </summary>
const ERROR_NOT_A_REPARSE_POINT = 4390;

/// <summary>
/// The reparse point attribute cannot be set because it conflicts with an existing attribute.
/// </summary>
const ERROR_REPARSE_ATTRIBUTE_CONFLICT = 4391;

/// <summary>
/// The data present in the reparse point buffer is invalid.
/// </summary>
const ERROR_INVALID_REPARSE_DATA = 4392;

/// <summary>
/// The tag present in the reparse point buffer is invalid.
/// </summary>
const ERROR_REPARSE_TAG_INVALID = 4393;

/// <summary>
/// There is a mismatch between the tag specified in the request and the tag present in the reparse point.
/// </summary>
const ERROR_REPARSE_TAG_MISMATCH = 4394;

/// <summary>
/// Command to set the reparse point data block.
/// </summary>
const FSCTL_SET_REPARSE_POINT = 0x000900A4;

/// <summary>
/// Command to get the reparse point data block.
/// </summary>
const FSCTL_GET_REPARSE_POINT = 0x000900A8;

/// <summary>
/// Command to delete the reparse point data base.
/// </summary>
const FSCTL_DELETE_REPARSE_POINT = 0x000900AC;

/// <summary>
/// Reparse point tag used to identify mount points and junction points.
/// </summary>
const IO_REPARSE_TAG_MOUNT_POINT = 0xA0000003;

/// <summary>
/// This prefix indicates to NTFS that the path is to be treated as a non-interpreted
/// path in the virtual file system.
/// </summary>
const NonInterpretedPathPrefix = `\??\`;

type EFileAccess  uint

var (
	EFileAccess_GenericRead EFileAccess = 0x80000000
	EFileAccess_GenericWrite EFileAccess = 0x40000000
	EFileAccess_GenericExecute EFileAccess = 0x20000000
	EFileAccess_GenericAll EFileAccess = 0x10000000
)

type  EFileShare uint

var (
	EFileShare_None EFileShare = 0x00000000
	EFileShare_Read EFileShare = 0x00000001
	EFileShare_Write EFileShare = 0x00000002
	EFileShare_Delete EFileShare = 0x00000004
)

type  ECreationDisposition  uint

var (
	ECreationDisposition_New ECreationDisposition = 1
	ECreationDisposition_CreateAlways ECreationDisposition = 2
	ECreationDisposition_OpenExisting ECreationDisposition = 3
	ECreationDisposition_OpenAlways ECreationDisposition = 4
	ECreationDisposition_TruncateExisting ECreationDisposition = 5
)

type  EFileAttributes  uint

var (
	EFileAttributes_Readonly EFileAttributes = 0x00000001
	EFileAttributes_Hidden EFileAttributes = 0x00000002
	EFileAttributes_System EFileAttributes = 0x00000004
	EFileAttributes_Directory EFileAttributes = 0x00000010
	EFileAttributes_Archive EFileAttributes = 0x00000020
	EFileAttributes_Device EFileAttributes = 0x00000040
	EFileAttributes_Normal EFileAttributes = 0x00000080
	EFileAttributes_Temporary EFileAttributes = 0x00000100
	EFileAttributes_SparseFile EFileAttributes = 0x00000200
	EFileAttributes_ReparsePoint EFileAttributes = 0x00000400
	EFileAttributes_Compressed EFileAttributes = 0x00000800
	EFileAttributes_Offline EFileAttributes = 0x00001000
	EFileAttributes_NotContentIndexed EFileAttributes = 0x00002000
	EFileAttributes_Encrypted EFileAttributes = 0x00004000
	EFileAttributes_Write_Through EFileAttributes = 0x80000000
	EFileAttributes_Overlapped EFileAttributes = 0x40000000
	EFileAttributes_NoBuffering EFileAttributes = 0x20000000
	EFileAttributes_RandomAccess EFileAttributes = 0x10000000
	EFileAttributes_SequentialScan EFileAttributes = 0x08000000
	EFileAttributes_DeleteOnClose EFileAttributes = 0x04000000
	EFileAttributes_BackupSemantics EFileAttributes = 0x02000000
	EFileAttributes_PosixSemantics EFileAttributes = 0x01000000
	EFileAttributes_OpenReparsePoint EFileAttributes = 0x00200000
	EFileAttributes_OpenNoRecall EFileAttributes = 0x00100000
	EFileAttributes_FirstPipeInstance EFileAttributes = 0x00080000
)

type   REPARSE_DATA_BUFFER struct {
	/// <summary>
	/// Reparse point tag. Must be a Microsoft reparse point tag.
	/// </summary>
	ReparseTag           uint

	/// <summary>
	/// Size, in bytes, of the data after the Reserved member. This can be calculated by:
	/// (4 * sizeof(ushort)) + SubstituteNameLength + PrintNameLength +
	/// (namesAreNullTerminated ? 2 * sizeof(char) : 0);
	/// </summary>
	ReparseDataLength    uint16

	/// <summary>
	/// Reserved; do not use.
	/// </summary>
	Reserved             uint16

	/// <summary>
	/// Offset, in bytes, of the substitute name string in the PathBuffer array.
	/// </summary>
	SubstituteNameOffset uint16

	/// <summary>
	/// Length, in bytes, of the substitute name string. If this string is null-terminated,
	/// SubstituteNameLength does not include space for the null character.
	/// </summary>
	SubstituteNameLength

	/// <summary>
	/// Offset, in bytes, of the print name string in the PathBuffer array.
	/// </summary>
	PrintNameOffset      uint16

	/// <summary>
	/// Length, in bytes, of the print name string. If this string is null-terminated,
	/// PrintNameLength does not include space for the null character.
	/// </summary>
	PrintNameLength      uint16

	/// <summary>
	/// A buffer containing the unicode-encoded path string. The path string contains
	/// the substitute name string and print name string.
	/// </summary>
	PathBuffer           []uint8
}

/*
syscall.DeviceIoControl(
handle Handle,
ioControlCode uint32,
inBuffer *byte,
inBufferSize uint32,
outBuffer *byte,
outBufferSize uint32,
bytesReturned *uint32,
overlapped *Overlapped) (err error) {
}
*/
func DeviceIoControl(handle Handle, ioControlCode uint32, inBuffer *byte, inBufferSize uint32, outBuffer *byte, outBufferSize uint32, bytesReturned *uint32, overlapped *Overlapped) (err error)  {
	return syscall.DeviceIoControl(handle,ioControlCode,inBuffer,inBufferSize,outBuffer,outBufferSize,bytesReturned,overlapped);
}
/*
[DllImport("kernel32.dll", CharSet = CharSet.Auto, SetLastError = true)]
private static extern bool DeviceIoControl(IntPtr hDevice, uint dwIoControlCode,
IntPtr InBuffer, int nInBufferSize,
IntPtr OutBuffer, int nOutBufferSize,
out int pBytesReturned, IntPtr lpOverlapped);
*/

func CreateFile(name *uint16, access uint32, mode uint32, sa *SecurityAttributes, createmode uint32, attrs uint32, templatefile int32) (handle Handle, err error) {
	return syscall.CreateFile(name,access,mode,sa,createmode,attrs,templatefile);
}
/*
[DllImport("kernel32.dll", SetLastError = true)]
private static extern IntPtr CreateFile(
string lpFileName,
EFileAccess dwDesiredAccess,
EFileShare dwShareMode,
IntPtr lpSecurityAttributes,
ECreationDisposition dwCreationDisposition,
EFileAttributes dwFlagsAndAttributes,
IntPtr hTemplateFile);
*/
/// <summary>
/// Creates a junction point from the specified directory to the specified target directory.
/// </summary>
/// <remarks>
/// Only works on NTFS.
/// </remarks>
/// <param name="junctionPoint">The junction point path</param>
/// <param name="targetDir">The target directory</param>
/// <param name="overwrite">If true overwrites an existing reparse point or empty directory</param>
/// <exception cref="IOException">Thrown when the junction point could not be created or when
/// an existing directory was found and <paramref name="overwrite" /> if false</exception>
func Create(junctionPoint string, targetDir string, overwrite bool) bool {
	targetDir,_ = filepath.Abs(targetDir)

	var exists bool
	if exists, _ = directory.DirectoryExists(targetDir); exists == false {
		fmt.Println(targetDir + " does not exist or is not a directory.")
		return false
		//log(targetDir + " does not exist or is not a directory.");
		//throw
		//new
		//IOException("Target path does not exist or is not a directory.");
		//return false;
	}

	if exists, _ = directory.DirectoryExists(junctionPoint); exists == false {
		if !overwrite {
			fmt.Println(junctionPoint + " already exists and overwrite parameter is false.")
			//log(junctionPoint + " already exists and overwrite parameter is false.");
			//throw
			//new
			//IOException("Directory already exists and overwrite parameter is false.");
			return false;
		}
	} else {
		os.MkdirAll(junctionPoint, 0777);

		//Directory.CreateDirectory(junctionPoint);
	}

	using(SafeFileHandlehandle = OpenReparsePoint(junctionPoint, EFileAccess.GenericWrite))
	{
		byte[] targetDirBytes = Encoding.Unicode.GetBytes(NonInterpretedPathPrefix + Path.GetFullPath(targetDir));

		REPARSE_DATA_BUFFER reparseDataBuffer = new REPARSE_DATA_BUFFER();

		reparseDataBuffer.ReparseTag = IO_REPARSE_TAG_MOUNT_POINT;
		reparseDataBuffer.ReparseDataLength = (ushort)(targetDirBytes.Length + 12);
		reparseDataBuffer.SubstituteNameOffset = 0;
		reparseDataBuffer.SubstituteNameLength = (ushort)targetDirBytes.Length;
		reparseDataBuffer.PrintNameOffset = (ushort)(targetDirBytes.Length + 2);
		reparseDataBuffer.PrintNameLength = 0;
		reparseDataBuffer.PathBuffer = new byte[0x3ff0];
		Array.Copy(targetDirBytes, reparseDataBuffer.PathBuffer, targetDirBytes.Length);

		int inBufferSize = Marshal.SizeOf(reparseDataBuffer);
		IntPtr inBuffer = Marshal.AllocHGlobal(inBufferSize);

		try
		{
			Marshal.StructureToPtr(reparseDataBuffer, inBuffer, false);

			int bytesReturned;
			bool result = DeviceIoControl(handle.DangerousGetHandle(), FSCTL_SET_REPARSE_POINT,
			inBuffer, targetDirBytes.Length + 20, IntPtr.Zero, 0, out bytesReturned, IntPtr.Zero);

			if (!result)
				ThrowLastWin32Error("Unable to create junction point.");
		}
		finally
		{
			Marshal.FreeHGlobal(inBuffer);
		}
	}
}

/// <summary>
/// Deletes a junction point at the specified source directory along with the directory itself.
/// Does nothing if the junction point does not exist.
/// </summary>
/// <remarks>
/// Only works on NTFS.
/// </remarks>
/// <param name="junctionPoint">The junction point path</param>
public static void Delete(string junctionPoint)
{
	if (!Directory.Exists(junctionPoint))
	{
		if (File.Exists(junctionPoint))
			throw new IOException("Path is not a junction point.");
		return;
	}

	using (SafeFileHandle handle = OpenReparsePoint(junctionPoint, EFileAccess.GenericWrite))
	{
		REPARSE_DATA_BUFFER reparseDataBuffer = new REPARSE_DATA_BUFFER();

		reparseDataBuffer.ReparseTag = IO_REPARSE_TAG_MOUNT_POINT;
		reparseDataBuffer.ReparseDataLength = 0;
		reparseDataBuffer.PathBuffer = new byte[0x3ff0];

		int inBufferSize = Marshal.SizeOf(reparseDataBuffer);
		IntPtr inBuffer = Marshal.AllocHGlobal(inBufferSize);
		try
		{
			Marshal.StructureToPtr(reparseDataBuffer, inBuffer, false);

			int bytesReturned;
			bool result = DeviceIoControl(handle.DangerousGetHandle(), FSCTL_DELETE_REPARSE_POINT,
				inBuffer, 8, IntPtr.Zero, 0, out bytesReturned, IntPtr.Zero);

			if (!result)
				ThrowLastWin32Error("Unable to delete junction point.");
		}
		finally
		{
			Marshal.FreeHGlobal(inBuffer);
		}

		try
		{
			Directory.Delete(junctionPoint);
		}
		catch (IOException ex)
		{
			throw new IOException("Unable to delete junction point.", ex);
		}
	}
}

/// <summary>
/// Determines whether the specified path exists and refers to a junction point.
/// </summary>
/// <param name="path">The junction point path</param>
/// <returns>True if the specified path represents a junction point</returns>
/// <exception cref="IOException">Thrown if the specified path is invalid
/// or some other error occurs</exception>
public static bool Exists(string path)
{
	if (! Directory.Exists(path))
		return false;
	using (SafeFileHandle handle = OpenReparsePoint(path, EFileAccess.GenericRead))
	{
		string target = InternalGetTarget(handle);
		return target != null;
	}
}

/// <summary>
/// Gets the target of the specified junction point.
/// </summary>
/// <remarks>
/// Only works on NTFS.
/// </remarks>
/// <param name="junctionPoint">The junction point path</param>
/// <returns>The target of the junction point</returns>
/// <exception cref="IOException">Thrown when the specified path does not
/// exist, is invalid, is not a junction point, or some other error occurs</exception>
public static string GetTarget(string junctionPoint)
{
	using (SafeFileHandle handle = OpenReparsePoint(junctionPoint, EFileAccess.GenericRead))
	{
		string target = InternalGetTarget(handle);
		if (target == null)
			throw new IOException("Path is not a junction point.");
		return target;
	}
}

func  InternalGetTarget(handle Handle) string
{
	var reparseDataBuffer REPARSE_DATA_BUFFER
	var outBufferSize int= unsafe.Sizeof(reparseDataBuffer)


	var outBuffer *[]byte= &[outBufferSize]byte{}

	try
	{
		var  bytesReturned *uint32;
		bool result = DeviceIoControl(handle, FSCTL_GET_REPARSE_POINT,
			&[0]byte{}, 0, outBuffer, outBufferSize, bytesReturned, &[0]byte{});

		if !result
		{

			if lasterr:=syscall.GetLastError(); lasterr== ERROR_NOT_A_REPARSE_POINT{
				fmt.print("Unable to open reparse point.")
				return nil
			}
			panic("Unable to open reparse point.")

		}

		REPARSE_DATA_BUFFER reparseDataBuffer = (REPARSE_DATA_BUFFER)Marshal.PtrToStructure(outBuffer, typeof(REPARSE_DATA_BUFFER));

		if (reparseDataBuffer.ReparseTag != IO_REPARSE_TAG_MOUNT_POINT)
			return null;

		string targetDir = Encoding.Unicode.GetString(reparseDataBuffer.PathBuffer,reparseDataBuffer.SubstituteNameOffset, reparseDataBuffer.SubstituteNameLength);

		if (targetDir.StartsWith(NonInterpretedPathPrefix))
			targetDir = targetDir.Substring(NonInterpretedPathPrefix.Length);

		return targetDir;
	}
	finally
	{
		Marshal.FreeHGlobal(outBuffer);
	}
}

func OpenReparsePoint(reparsePoint string, accessMode EFileAccess) Handle{
	handle , err :=CreateFile(reparsePoint,accessMode,EFileShare_Read | EFileShare_Write | EFileShare_Delete,&SecurityAttributes{Length:uint32(unsafe.Sizeof(sa)),InheritHandle:0},	ECreationDisposition_OpenExisting,EFileAttributes_BackupSemantics | EFileAttributes_OpenReparsePoint,0)

	//SafeFileHandle reparsePointHandle = new SafeFileHandle(CreateFile(reparsePoint, accessMode,EFileShare.Read | EFileShare.Write | EFileShare.Delete,
	//	IntPtr.Zero, ECreationDisposition.OpenExisting,EFileAttributes.BackupSemantics | EFileAttributes.OpenReparsePoint, IntPtr.Zero), true);

	if lasterr:=syscall.GetLastError(); lasterr!=null{
		panic("Unable to open reparse point.")
	}
	//if (Marshal.GetLastWin32Error() != 0)
	//	ThrowLastWin32Error("Unable to open reparse point.");

	return handle
}

//private static void ThrowLastWin32Error(string message)
//{
//	throw new IOException(message, Marshal.GetExceptionForHR(Marshal.GetHRForLastWin32Error()));
//}
}
}
