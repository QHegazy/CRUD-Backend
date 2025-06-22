package my_utils

const key = 14
const offset = 78
const prime = 21

func ObfuscateNumbers(number uint64) uint64 {

	transformed := (number ^ key) * prime

	scrambled := ((transformed + offset) << 7) | ((transformed + offset) >> (64 - 7))

	return scrambled ^ (key * prime)
}

func DeobfuscateNumbers(obfuscatedNumber uint64) uint64 {

	unscrambled := obfuscatedNumber ^ (key * prime)

	reversed := (unscrambled >> 7) | (unscrambled << (64 - 7))

	original := ((reversed - offset) / prime) ^ key

	return original
}
