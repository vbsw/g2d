#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WIN32)

typedef unsigned long g2d_ul_t;
extern void g2d_free(void *data);
extern void g2d_init(int *err_num, g2d_ul_t *err_win32);

#elif defined(G2D_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_H */