#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_MODULES_WIN32)

/*
typedef unsigned long g2d_ul_t;
typedef unsigned int g2d_ui_t;
extern void g2d_mods_rects_init(void **data, int *err_num, g2d_ul_t *err_win32, char **err_str);
*/

typedef struct { void **all; char *err_str; void *set_func, *get_func; long long err1, err2; int list_len, list_cap, words_len, words_cap; } cdata_t;
extern void g2d_mods_rects_init(int pass, cdata_t *cdata);

#elif defined(G2D_MODULES_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_H */