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
	GenericRead EFileAccess = 0x80000000
	GenericWrite EFileAccess = 0x40000000
	GenericExecute EFileAccess = 0x20000000
	GenericAll EFileAccess = 0x10000000
)

type  EFileShare uint

var (
	None EFileShare = 0x00000000
	Read EFileShare = 0x00000001
	Write EFileShare = 0x00000002
	Delete EFileShare = 0x00000004
)

type  ECreationDisposition  uint

var (
	New ECreationDisposition = 1
	CreateAlways ECreationDisposition = 2
	OpenExisting ECreationDisposition = 3
	OpenAlways ECreationDisposition = 4
	TruncateExisting ECreationDisposition = 5
)

type  EFileAttributes  uint

var (
	Readonly EFileAttributes = 0x00000001
	Hidden EFileAttributes = 0x00000002
	System EFileAttributes = 0x00000004
	Directory EFileAttributes = 0x00000010
	Archive EFileAttributes = 0x00000020
	Device EFileAttributes = 0x00000040
	Normal EFileAttributes = 0x00000080
	Temporary EFileAttributes = 0x00000100
	SparseFile EFileAttributes = 0x00000200
	ReparsePoint EFileAttributes = 0x00000400
	Compressed EFileAttributes = 0x00000800
	Offline EFileAttributes = 0x00001000
	NotContentIndexed EFileAttributes = 0x00002000
	Encrypted EFileAttributes = 0x00004000
	Write_Through EFileAttributes = 0x80000000
	Overlapped EFileAttributes = 0x40000000
	NoBuffering EFileAttributes = 0x20000000
	RandomAccess EFileAttributes = 0x10000000
	SequentialScan EFileAttributes = 0x08000000
	DeleteOnClose EFileAttributes = 0x04000000
	BackupSemantics EFileAttributes = 0x02000000
	PosixSemantics EFileAttributes = 0x01000000
	OpenReparsePoint EFileAttributes = 0x00200000
	OpenNoRecall EFileAttributes = 0x00100000
	FirstPipeInstance EFileAttributes = 0x00080000
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

[DllImport("kernel32.dll", CharSet = CharSet.Auto, SetLastError = true)]
private static extern bool DeviceIoControl(IntPtr hDevice, uint dwIoControlCode,
IntPtr InBuffer, int nInBufferSize,
IntPtr OutBuffer, int nOutBufferSize,
out int pBytesReturned, IntPtr lpOverlapped);

[DllImport("kernel32.dll", SetLastError = true)]
private static extern IntPtr CreateFile(
string lpFileName,
EFileAccess dwDesiredAccess,
EFileShare dwShareMode,
IntPtr lpSecurityAttributes,
ECreationDisposition dwCreationDisposition,
EFileAttributes dwFlagsAndAttributes,
IntPtr hTemplateFile);

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
public static bool Create(string junctionPoint, string targetDir, bool overwrite)
{
targetDir = Path.GetFullPath(targetDir);

if (!Directory.Exists(targetDir))
{
//log(targetDir + " does not exist or is not a directory.");
throw new IOException("Target path does not exist or is not a directory.");
//return false;
}

if (Directory.Exists(junctionPoint))
{
if (!overwrite) {
//log(junctionPoint + " already exists and overwrite parameter is false.");
throw new IOException("Directory already exists and overwrite parameter is false.");
//return false;
}
}
else
{
Directory.CreateDirectory(junctionPoint);
}

using (TTuple<bool, SafeFileHandle> ret = OpenReparsePoint(junctionPoint, EFileAccess.GenericWrite))
{
if (!ret.result){
return false;
}
SafeFileHandle handle = ret.data;
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

if (!result) {
//log("Unable to create junction point : "+junctionPoint+" to : "+targetDir);
ThrowLastWin32Error("Unable to create junction point.");
//return false;
}
return true;
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
public static bool Delete(string junctionPoint)
{
if (!Directory.Exists(junctionPoint))
{
if (File.Exists(junctionPoint)) {
//log(junctionPoint + " is not a junction point.");
throw new IOException("Path is not a junction point.");
//return false;
}

}

using (TTuple<bool, SafeFileHandle> ret = OpenReparsePoint(junctionPoint, EFileAccess.GenericWrite))
{
if (!ret.result){
return false;
}
SafeFileHandle handle = ret.data;

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
{
ThrowLastWin32Error("Unable to delete junction point.");
//log("Unable to delete junction point :" + junctionPoint);
//return false;
}

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
//log("Unable to delete junction point :" + junctionPoint+ ","+ex.ToString());
throw new IOException("Unable to delete junction point.", ex);
}
return true;
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

using (TTuple<bool, SafeFileHandle> ret = OpenReparsePoint(path, EFileAccess.GenericRead))
{
if (!ret.result)
{
return false;
}
SafeFileHandle handle = ret.data;
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
using (TTuple<bool, SafeFileHandle> ret = OpenReparsePoint(junctionPoint, EFileAccess.GenericRead))
{
if (!ret.result)
{
//不会到这里
return "";
}
SafeFileHandle handle = ret.data;
string target = InternalGetTarget(handle);
if (target == null)
{
throw new IOException("Path is not a junction point.");
//log("Path is not a junction point.");
//return "";

}
return target;
}
}

private static string InternalGetTarget(SafeFileHandle handle)
{
int outBufferSize = Marshal.SizeOf(typeof(REPARSE_DATA_BUFFER));
IntPtr outBuffer = Marshal.AllocHGlobal(outBufferSize);

try
{
int bytesReturned;
bool result = DeviceIoControl(handle.DangerousGetHandle(), FSCTL_GET_REPARSE_POINT,
IntPtr.Zero, 0, outBuffer, outBufferSize, out bytesReturned, IntPtr.Zero);

if (!result)
{
int error = Marshal.GetLastWin32Error();
if (error == ERROR_NOT_A_REPARSE_POINT)
return null;

ThrowLastWin32Error("Unable to get information about junction point.");
}

REPARSE_DATA_BUFFER reparseDataBuffer = (REPARSE_DATA_BUFFER)
Marshal.PtrToStructure(outBuffer, typeof(REPARSE_DATA_BUFFER));

if (reparseDataBuffer.ReparseTag != IO_REPARSE_TAG_MOUNT_POINT)
return null;

string targetDir = Encoding.Unicode.GetString(reparseDataBuffer.PathBuffer,
reparseDataBuffer.SubstituteNameOffset, reparseDataBuffer.SubstituteNameLength);

if (targetDir.StartsWith(NonInterpretedPathPrefix))
targetDir = targetDir.Substring(NonInterpretedPathPrefix.Length);

return targetDir;
}
finally
{
Marshal.FreeHGlobal(outBuffer);
}
}

private static TTuple<bool, SafeFileHandle> OpenReparsePoint(string reparsePoint, EFileAccess accessMode)
{
SafeFileHandle reparsePointHandle = new SafeFileHandle(CreateFile(reparsePoint, accessMode,
EFileShare.Read | EFileShare.Write | EFileShare.Delete,
IntPtr.Zero, ECreationDisposition.OpenExisting,
EFileAttributes.BackupSemantics | EFileAttributes.OpenReparsePoint, IntPtr.Zero), true);

if (Marshal.GetLastWin32Error() != 0)
{
//log("Unable to open reparse point.");
//return TTuple<bool,SafeFileHandle>.create(false, reparsePointHandle);
ThrowLastWin32Error("Unable to open reparse point.");
}

//return reparsePointHandle;
return TTuple<bool, SafeFileHandle>.create(true, reparsePointHandle);
}

private static void ThrowLastWin32Error(string message)
{
throw new IOException(message, Marshal.GetExceptionForHR(Marshal.GetHRForLastWin32Error()));
}

private static void log(string log) {
StreamWriter w = File.AppendText("JunctionPoint.log");
w.WriteLine("{0}:{1}", System.DateTime.Now, log);
w.Close();
}

public class TTuple<T1, T2> : System.IDisposable
{
public readonly T1 result;
public readonly T2 data;
public void Dispose()
{
}

public TTuple(T1 result, T2 data)
{
this.result = result;
this.data = data;
}

public static TTuple<T1, T2> create<T1, T2>(T1 result, T2 data)
{
return new TTuple<T1, T2>(result, data);
}
}

//static void Main(string[] args)
//{
//   Console.WriteLine(JunctionPoint.Create("v:/ttt", "v:/temp", true));
//  Console.WriteLine(JunctionPoint.GetTarget(@"v:/ttt"));
// Console.WriteLine(JunctionPoint.Exists(@"v:/ttt"));
//  Console.WriteLine(JunctionPoint.Delete(@"v:/ttt"));
//}
