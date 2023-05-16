package rabin_karp

const Base = 256
const Mod = 997

func RabinKarpSearch(s, pattern string) (int, bool) {
	M := len(pattern)
	N := len(s)
	if M > N {
		return 0, false
	}
	H := 1
	for i := 0; i < M-1; i++ {
		H = (H * Base) % Mod
	}

	h0 := initHash(pattern)
	h1 := initHash(s[0:M])
	for i := 0; i <= N-M; i++ {
		if h0 == h1 && s[i:i+M] == pattern {
			return i, true
		}

		// if not reach the end of s
		if i < N-M {
			h1 = rollingHash(s, i, i+M, H, h1)
		}
	}
	return 0, false
}
func initHash(str string) int {
	h := 0
	for _, c := range str {
		h = h*Base + int(c)
		h = h % Mod
	}
	return h
}

// next rolling hash
func rollingHash(str string, i, j, H, prev int) int {
	// trim the leading digit and add the trailing digit
	h := ((prev-int(str[i])*H)*Base + int(str[j])) % Mod
	if h < 0 {
		h += Mod
	}
	return h
}
