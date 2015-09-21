#include <stdio.h>
#include <assert.h>
#include <stdlib.h>

#include "dh64.h"

// cc dh64.c dh64_test.c
int main(int argc, char **argv) {
	int n = 10000;
	if (argc > 1) {
		n = atoi(argv[1]);
	}
	for (int i = 0; i < n; i ++) {
		uint64_t private_key1;
		uint64_t public_key1;
		dh64_key_pair(&private_key1, &public_key1);


		uint64_t private_key2;
		uint64_t public_key2;
		dh64_key_pair(&private_key2, &public_key2);
		
		uint64_t secret1 = dh64_secret(private_key1, public_key2);
		uint64_t secret2 = dh64_secret(private_key2, public_key1);

		assert(secret1 == secret2);
		printf("{0x%016llX, 0x%016llX, 0x%016llX},\n", public_key1, private_key1, secret1);
		printf("{0x%016llX, 0x%016llX, 0x%016llX},\n", public_key2, private_key2, secret2);
	}
}
