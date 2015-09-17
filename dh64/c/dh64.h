#ifndef DH64_H
#define DH64_H

#include <stdint.h>

typedef struct dh64 {
	uint64_t public_key;
	uint64_t private_key;
	uint64_t secret;
} dh64;

dh64 dh64_gen();

void dh64_init(dh64* dh, const uint64_t another_public);

#endif
