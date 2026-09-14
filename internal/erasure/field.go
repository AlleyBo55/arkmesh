package erasure

// Arithmetic in GF(2^8) using the reducing polynomial x^8 + x^4 + x^3 + x^2 + 1
// (0x11d). The generator 2 is primitive for this polynomial, which a test proves
// by requiring every nonzero element to appear exactly once in the exponent
// table. Reed Solomon reconstruction silently fails if that property does not
// hold, so it is checked rather than assumed.

const (
	fieldElements  = 256
	reduce         = 0x1d
	fieldGenerator = 2
)

var (
	expTable [512]byte
	logTable [fieldElements]byte
)

func init() {
	value := byte(1)
	for power := 0; power < fieldElements-1; power++ {
		expTable[power] = value
		logTable[value] = byte(power)
		value = mulGenerator(value)
	}
	for power := fieldElements - 1; power < len(expTable); power++ {
		expTable[power] = expTable[power-(fieldElements-1)]
	}
}

// mulGenerator multiplies by the generator, reducing modulo the field polynomial.
func mulGenerator(value byte) byte {
	high := value & 0x80
	shifted := value << 1
	if high != 0 {
		shifted ^= reduce
	}
	return shifted
}

func mul(left, right byte) byte {
	if left == 0 || right == 0 {
		return 0
	}
	return expTable[int(logTable[left])+int(logTable[right])]
}

func div(numerator, denominator byte) byte {
	if denominator == 0 {
		panic("erasure: division by zero in GF(2^8)")
	}
	if numerator == 0 {
		return 0
	}
	return expTable[int(logTable[numerator])-int(logTable[denominator])+(fieldElements-1)]
}

func inverse(value byte) byte {
	if value == 0 {
		panic("erasure: zero has no inverse in GF(2^8)")
	}
	return expTable[(fieldElements-1)-int(logTable[value])]
}

func exponent(base byte, power int) byte {
	if power == 0 {
		return 1
	}
	if base == 0 {
		return 0
	}
	return expTable[(int(logTable[base])*power)%(fieldElements-1)]
}
