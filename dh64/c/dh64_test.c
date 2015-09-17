#include <stdio.h>
#include <assert.h>
#include <stdlib.h>

#include "dh64.h"

// cc dh64.c dh64_test.c
int main(void) {
	for (int i = 0; i < 10000; i ++) {
		dh64 dh1 = dh64_gen();
		dh64 dh2 = dh64_gen();

		dh64_init(&dh1, dh2.public_key);
		dh64_init(&dh2, dh1.public_key);

		assert(dh1.secret == dh2.secret);

		printf("dh1 = %16llx, %16llx, %16llx\n", dh1.public_key, dh1.private_key, dh1.secret);
		printf("dh2 = %16llx, %16llx, %16llx\n", dh2.public_key, dh2.private_key, dh2.secret);
	}
}
