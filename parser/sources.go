package parser

import "fmt"

var (
	sourceMap map[string]string
	aeadNames []string
)

func init() {
	sourceMap = map[string]string{
		"GOES": "Geosynchronous Orbit Environment Satellite",
		"GPS":  "Global Position System",
		"GAL":  "Galileo Positioning System",
		"PPS":  "Generic pulse-per-second",
		"IRIG": "Inter-Range Instrumentation Group",
		"WWVB": "LF Radio WWVB Ft. Collins, CO 60 kHz",
		"DCF":  "LF Radio DCF77 Mainflingen, DE 77.5 kHz",
		"HBG":  "LF Radio HBG Prangins, HB 75 kHz",
		"MSF":  "LF Radio MSF Anthorn, UK 60 kHz",
		"JJY":  "LF Radio JJY Fukushima, JP 40 kHz, Saga, JP 60 kHz",
		"LORC": "MF Radio LORAN C station, 100 kHz",
		"TDF":  "MF Radio Allouis, FR 162 kHz",
		"CHU":  "HF Radio CHU Ottawa, Ontario",
		"WWV":  "HF Radio WWV Ft. Collins, CO",
		"WWVH": "HF Radio WWVH Kauai, HI",
		"NIST": "NIST telephone modem",
		"ACTS": "NIST telephone modem",
		"USNO": "USNO telephone modem",
		"PTB":  "European telephone modem",
	}
	aeadNames = []string{
		"",
		"AEAD_AES_128_GCM",           // 1
		"AEAD_AES_256_GCM",           // 2
		"AEAD_AES_128_CCM",           // 3
		"AEAD_AES_256_CCM",           // 4
		"AEAD_AES_128_GCM_8",         // 5
		"AEAD_AES_256_GCM_8",         // 6
		"AEAD_AES_128_GCM_12",        // 7
		"AEAD_AES_256_GCM_12",        // 8
		"AEAD_AES_128_CCM_SHORT",     // 9
		"AEAD_AES_256_CCM_SHORT",     // 10
		"AEAD_AES_128_CCM_SHORT_8",   // 11
		"AEAD_AES_256_CCM_SHORT_8",   // 12
		"AEAD_AES_128_CCM_SHORT_12",  // 13
		"AEAD_AES_256_CCM_SHORT_12",  // 14
		"AEAD_AES_SIV_CMAC_256",      // 15
		"AEAD_AES_SIV_CMAC_384",      // 16
		"AEAD_AES_SIV_CMAC_512",      // 17
		"AEAD_AES_128_CCM_8",         // 18
		"AEAD_AES_256_CCM_8",         // 19
		"AEAD_AES_128_OCB_TAGLEN128", // 20
		"AEAD_AES_128_OCB_TAGLEN96",  // 21
		"AEAD_AES_128_OCB_TAGLEN64",  // 22
		"AEAD_AES_192_OCB_TAGLEN128", // 23
		"AEAD_AES_192_OCB_TAGLEN96",  // 24
		"AEAD_AES_192_OCB_TAGLEN64",  // 25
		"AEAD_AES_256_OCB_TAGLEN128", // 26
		"AEAD_AES_256_OCB_TAGLEN96",  // 27
		"AEAD_AES_256_OCB_TAGLEN64",  // 28
		"AEAD_CHACHA20_POLY1305",     // 29
		"AEAD_AES_128_GCM_SIV",       // 30
		"AEAD_AES_256_GCM_SIV",       // 31
		"AEAD_AEGIS128L",             // 32
		"AEAD_AEGIS256",              // 33
	}
}

func completeSource(s []byte) string {
	str := string(s)
	if complete, ok := sourceMap[str]; ok {
		return fmt.Sprintf("%s (%s)", str, complete)
	}
	if complete, ok := sourceMap[str[:3]]; ok {
		return fmt.Sprintf("%s (%s)", str, complete)
	}
	return str
}

func getAEADName(id byte) string {
	return aeadNames[id]
}
