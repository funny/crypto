// The biggest 64bit prime
#define P 0xffffffffffffffc5ull
#define G 5

#include <stdlib.h>

#include "dh64.h"

// calc a * b % p , avoid 64bit overflow
static inline uint64_t
mul_mod_p(uint64_t a, uint64_t b) {
	uint64_t m = 0;
	while (b) {
		if (b&1) {
			uint64_t t = P - a;
			if (m >= t) {
				m -= t;
			} else {
				m += a;
			}
		}
		if (a >= P - a) {
			a = a * 2 - P;
		} else {
			a = a * 2;
		}
		b >>= 1;
	}
	return m;
}

static inline uint64_t
pow_mod_p(uint64_t a, uint64_t b) {
	if (b == 1) {
		return a;
	}
	uint64_t t = pow_mod_p(a, b>>1);
	t = mul_mod_p(t, t);
	if (b % 2) {
		t = mul_mod_p(t, a);
	}
	return t;
}

// calc a^b % p
static inline uint64_t
powmodp(uint64_t a, uint64_t b) {
	if (a == 0)
		return 1;
	if (b == 0)
		return 1;
	if (a > P)
		a %= P;
	return pow_mod_p(a, b);
}

void 
dh64_key_pair(uint64_t* private_key, uint64_t* public_key) {
	uint64_t a = rand();
	uint64_t b = rand() & 0xFFFF;
	uint64_t c = rand() & 0xFFFF;
	uint64_t d = (rand() & 0xFFFF) + 1;
	*private_key = a << 48 | b << 32 | c << 16 | d;
	*public_key  = powmodp(G, *private_key);
}

uint64_t 
dh64_public_key(const uint64_t private_key) {
	return powmodp(G, private_key);
}

uint64_t
dh64_secret(const uint64_t private_key, const uint64_t another_public_key) {
	return powmodp(another_public_key, private_key);
}
