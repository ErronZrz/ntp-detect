package experiment

var (
	commaSpace   = []byte{0x2C, 0x20}
	commaNewLine = []byte{0x2C, 0x0D, 0x0A}
)

/*
func main() {
	csStr := string(commaSpace)
	clStr := string(commaNewLine)
	r := []rune("haha")
	fmt.Println(unsafe.Sizeof(csStr), csStr)
	fmt.Println(unsafe.Sizeof(clStr), clStr)
	fmt.Println(unsafe.Sizeof(r), r)
	hello := "Hello, world!"
	fmt.Println(len(hello), unsafe.Sizeof(hello))

	longStr := "1234567890123456789012345678901234567890"
	fmt.Println(len(longStr), unsafe.Sizeof(longStr))

	testRune := []rune{'a', 'b', 'c'}
	runeStr := string(testRune)
	fmt.Println(runeStr)
	testRune[1] = 'd'
	fmt.Println(runeStr)
	sa := "hello"
	sb := "hello"
	sc := sa
	scc := sa[1:5]
	bd := make([]byte, 5, 100)
	bd[0], bd[1], bd[2], bd[3], bd[4] = 'h', 'e', 'l', 'l', 'o'
	sd := string(bd)
	se := string([]byte{'h', 'e', 'l', 'l', 'o'})
	sf := sd
	sff := sd[1:5]
	gr := make([]rune, 5, 100)
	gr[0], gr[1], gr[2], gr[3], gr[4] = 'h', 'e', 'l', 'l', 'o'
	sg := string(gr)
	fmt.Println(len(sg))
	sh := string([]rune{'h', 'e', 'l', 'l', 'o'})
	si := sg
	sii := sg[1:5]
	sap := *(*unsafe.Pointer)(unsafe.Pointer(&sa))
	sbp := *(*unsafe.Pointer)(unsafe.Pointer(&sb))
	scp := *(*unsafe.Pointer)(unsafe.Pointer(&sc))
	sccp := *(*unsafe.Pointer)(unsafe.Pointer(&scc))
	sdp := *(*unsafe.Pointer)(unsafe.Pointer(&sd))
	sep := *(*unsafe.Pointer)(unsafe.Pointer(&se))
	sfp := *(*unsafe.Pointer)(unsafe.Pointer(&sf))
	sffp := *(*unsafe.Pointer)(unsafe.Pointer(&sff))
	sgp := *(*unsafe.Pointer)(unsafe.Pointer(&sg))
	shp := *(*unsafe.Pointer)(unsafe.Pointer(&sh))
	sip := *(*unsafe.Pointer)(unsafe.Pointer(&si))
	siip := *(*unsafe.Pointer)(unsafe.Pointer(&sii))
	fmt.Println(sap, sbp, scp, sccp)
	fmt.Println(sdp, sep, sfp, sffp)
	fmt.Println(sgp, shp, sip, siip)
	fmt.Println(*(*[16]byte)(sap))
	fmt.Println(*(*[16]byte)(sdp))
	fmt.Println(*(*[16]byte)(sgp))
	fmt.Println(*(*[16]byte)(*(*unsafe.Pointer)(unsafe.Pointer(&[]rune{'h', 'e', 'l', 'l', 'o'}))))
}
*/
