#ifndef VBSW_G2D_WINDOW_H
#define VBSW_G2D_WINDOW_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WINDOW_WIN32)

#include <stdio.h>

typedef struct { void *set_func, *get_func; char *err_str; void **all; long long err1, err2; int list_len, list_cap, words_len, words_cap; } cdata_t;
extern void g2d_window_init(int pass, cdata_t *cdata);

extern void g2d_window_mainloop();
extern void g2d_window_post_custom_msg(long long *millis, long long *err1, long long *err2);
extern void g2d_window_post_quit_msg(long long *err1, long long *err2);
extern void g2d_window_mainloop_clean_up();

extern void g2d_window_create(void **data, int cb_id, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void *t, size_t ts, long long *err1, long long *err2);
extern void g2d_window_show(void *data, long long *err1, long long *err2);

#elif defined(G2D_WINDOW_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* VBSW_G2D_WINDOW_H */