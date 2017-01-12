#include "blosc.h"

#include "bitshuffle-generic.c"
#include "shuffle-generic.c"

#if defined(__AVX2__)
#include "bitshuffle-avx2.c"
#include "shuffle-avx2.c"
#endif

#if defined(__SSE2__)
#include "bitshuffle-sse2.c"
#include "shuffle-sse2.c"
#endif


#include "blosc.c"
#include "blosclz.c"
#include "shuffle.c"
