package cache

func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	// Faster loop for larger strings
	for l >= 4 {
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
		d = (d * 33) ^ uint32(k[i+3])
		i += 4
		l -= 4
	}
	// Handle remaining bytes
	for l > 0 {
		d = (d * 33) ^ uint32(k[i])
		i++
		l--
	}
	return d ^ (d >> 16)
}
