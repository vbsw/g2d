#ifndef VBSW_G2D_OPENGLLOADER_H
#define VBSW_G2D_OPENGLLOADER_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_OPENGLLOADER_WIN32)

typedef struct { void **all; char *err_str; void *set_func, *get_func; long long err1, err2; int list_len, list_cap, words_len, words_cap; } cdata_t;
extern void g2d_openglloader_init(int pass, cdata_t *cdata);

#elif defined(G2D_OPENGLLOADER_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* VBSW_G2D_OPENGLLOADER_H */