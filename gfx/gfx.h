#ifndef G2D_GFX_H
#define G2D_GFX_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_GFX_WIN32)

typedef struct { void *set_func, *get_func; char *err_str; void **all; long long err1, err2; int list_len, list_cap, words_len, words_cap; } cdata_t;
extern void g2d_gfx_rects_init(int pass, cdata_t *cdata);

#elif defined(G2D_GFX_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_GFX_H */