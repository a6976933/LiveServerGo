package av

const ( //Header Signature
	F = 0x46
	L = 0x4c
	V = 0x56
)
const ( //FLV Verstion
	FLV_VERSION = 0x01
)
const ( //TypeFlag
	VIDEO_ONLY      = 0x01
	AUDIO_ONLY      = 0x04
	VIDEO_AND_AUDIO = 0x05
)
const ( //Sound Format
	LINEAR_PCM_PLATFORM_ENDIAN = iota
	ADPCM
	MP3
	LINEAR_PCM_LITTLE_ENDIAN
	NELLYMOSER_16_KHZ_MONO
	NELLYMOSER_8_KHZ_MONO
	NELLYMOSER
	G711_A_LAW_LOGARITHMIC_PCM
	G711_MU_LAW_LOGARITHMIC_PCM
	RESERVED
	AAC
	SPEEX
	_
	_
	MP3_8_KHZ
	DEVICE_SPECIFIC_SOUND
)
const ( //Sound Rate
	SOUND_5_5KHZ = iota
	SOUND_11KHZ
	SOUND_22KHZ
	SOUND_44KHZ
)
const ( //Sound Size
	SOUND_8_BIT_SAMPLE = iota
	SOUND_16_BIT_SAMPLE
)
const ( //Sound Type
	MONO = iota
	STEREO
)
const ( //AAC Type
	AAC_SEQUENCE_HEADER = iota
	AAC_RAW
)
const ( //Frame Type
	KEY_FRAME = iota + 1
	INTER_FRAME
	DISPOSABLE_INTER_FRAME
	GENERATED_KEY_FRAME
	VIDEO_INFO_COMMAND_FRAME
)
const ( //CodecID
	SORENSON_H_263 = iota + 2
	SCREEN_VIDEO
	VP6
	VP6_ALPHA_CHANNEL
	SCREEN_VIDEO_VERSION_2
	AVC_H_264
)
const ( //AVC Packet Type
	AVC_SEQUENCE_HEADER = iota
	AVC_NALU
	AVC_EOS
)

//func DTSToSecond()
const ( //AVC Decoder Configuration Record
	CONFIGURATION_VERSION = 1
)

func Trans4Byte2UINT32(b []byte) uint32 {
	return uint32(b[0])*16777216 + uint32(b[1])*65536 + uint32(b[2])*256 + uint32(b[3])
}

func Trans3Byte2UINT32(b []byte) uint32 {
	return uint32(b[0])*65536 + uint32(b[1])*256 + uint32(b[2])
}

func Trans2Byte2UINT16(b []byte) uint16 {
	return uint16(b[0])*256 + uint16(b[1])
}

func Trans1Byte2UINT8(b []byte) uint8 {
	return uint8(b[0])
}

func TransUINT32_2_4Byte(i uint32) []byte {
	var b = make([]byte, 4)
	b[0] = uint8(i >> 24)
	b[1] = uint8((i >> 16) & 255)
	b[1] = uint8((i >> 8) & 255)
	b[2] = uint8(i & 255)
	return b
}

func TransUINT32_2_3Byte(i uint32) []byte {
	var b = make([]byte, 3)
	b[0] = uint8(i >> 16)
	b[1] = uint8((i >> 8) & 255)
	b[2] = uint8(i & 255)
	return b
}

func TransUINT8_2_1Byte(i uint8) []byte {
	var b = make([]byte, 1)
	b[0] = uint8(i)
	return b
}

func ShiftBytesRight(b []byte, i uint) []byte {
	b = b[i:]
	return b
}
