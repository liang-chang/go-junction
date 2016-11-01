package junction

//https://msdn.microsoft.com/en-us/library/cc232006.aspx
//https://gist.github.com/Perlmint/f9f0e37db163dd69317d
//https://github.com/golang/go/blob/master/src/os/os_windows_test.go

/// The file or directory is not a reparse point.
const ERROR_NOT_A_REPARSE_POINT = 4390;

/// The reparse point attribute cannot be set because it conflicts with an existing attribute.
const ERROR_REPARSE_ATTRIBUTE_CONFLICT = 4391;

/// The data present in the reparse point buffer is invalid.
const ERROR_INVALID_REPARSE_DATA = 4392;

/// The tag present in the reparse point buffer is invalid.
const ERROR_REPARSE_TAG_INVALID = 4393;

/// There is a mismatch between the tag specified in the request and the tag present in the reparse point.
const ERROR_REPARSE_TAG_MISMATCH = 4394;

/// Command to set the reparse point data block.
const FSCTL_SET_REPARSE_POINT = 0x000900A4;

/// Command to get the reparse point data block.
const FSCTL_GET_REPARSE_POINT = 0x000900A8;

/// Command to delete the reparse point data base.
const FSCTL_DELETE_REPARSE_POINT = 0x000900AC;

/// Reparse point tag used to identify mount points and junction points.
const IO_REPARSE_TAG_MOUNT_POINT = 0xA0000003;

/// This prefix indicates to NTFS that the path is to be treated as a non-interpreted
/// path in the virtual file system.
const NonInterpretedPathPrefix = `\??\`;

//https://msdn.microsoft.com/en-us/library/ff552012.aspx
//_REPARSE_DATA_BUFFER header
type _REPARSE_DATA_BUFFER struct {
	ReparseTag        uint32
	ReparseDataLength uint16
	Reserved          uint16
}

//https://msdn.microsoft.com/en-us/library/ff552012.aspx
//MountPointReparseBuffer
type   MountPointReparseBuffer struct {
	/// Reparse point tag. Must be a Microsoft reparse point tag.
	ReparseTag           uint32

	/// Size, in bytes, of the data after the Reserved member. This can be calculated by:
	/// (4 * sizeof(ushort)) + SubstituteNameLength + PrintNameLength +
	/// (namesAreNullTerminated ? 2 * sizeof(char) : 0);
	ReparseDataLength    uint16

	/// Reserved; do not use.
	Reserved             uint16

	/// Offset, in bytes, of the substitute name string in the PathBuffer array.
	SubstituteNameOffset uint16

	/// Length, in bytes, of the substitute name string. If this string is null-terminated,
	/// SubstituteNameLength does not include space for the null character.
	SubstituteNameLength uint16

	/// Offset, in bytes, of the print name string in the PathBuffer array.
	PrintNameOffset      uint16

	/// Length, in bytes, of the print name string. If this string is null-terminated,
	/// PrintNameLength does not include space for the null character.
	PrintNameLength      uint16

	/// A buffer containing the unicode-encoded path string. The path string contains
	/// the substitute name string and print name string.
	//PathBuffer           [1]uint16

	// SubstituteName - 264 widechars = 528 bytes
	// PrintName      - 260 widechars = 520 bytes
	//                                = 1048 bytes total
	PathBuffer           [1048]uint16
}




