#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WIN32)

#include <stdio.h>

extern void g2d_free(void *data);
extern void g2d_init(int *numbers, long long *err1, long long *err2, char **err_nfo);
extern void g2d_main_loop();
extern void g2d_post_request(long long *err1, long long *err2);
extern void g2d_post_quit(long long *err1, long long *err2);
extern void g2d_clean_up();
extern void g2d_window_create(void **data, int cb_id, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void *t, size_t ts, long long *err1, long long *err2);
extern void g2d_window_show(void *data, long long *err1, long long *err2);
extern void g2d_window_props(void *data, int *mx, int *my, int *x, int *y, int *w, int *h, int *wn, int *hn, int *wx, int *hx, int *b, int *d, int *r, int *f, int *l);
extern void g2d_window_destroy(void *data, long long *err1, long long *err2);

extern void g2d_window_pos_size_set(void *data, int x, int y, int width, int height);
extern void g2d_window_style_set(void *data, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l);
extern void g2d_window_fullscreen_set(void *data, long long *err1, long long *err2);
extern void g2d_window_restore_bak(void *data);
extern void g2d_window_pos_apply(void *data, long long *err1, long long *err2);
extern void g2d_window_move(void *data, long long *err1, long long *err2);
extern void g2d_window_title_set(void *data, void *t, size_t ts, long long *err1, long long *err2);
extern void g2d_mouse_pos_set(void *data, int x, int y, long long *err1, long long *err2);

extern void g2d_gfx_init(void *data, long long *err1, long long *err2, char **err_nfo);
extern void g2d_gfx_release(void *data, long long *err1, long long *err2);
extern void g2d_gfx_draw(void *data, int w, int h, int i, float r, float g, float b, float **buffs, const int *bs, void **procs, int l, long long *err1, long long *err2);
extern void g2d_gfx_draw_rectangles(void *data, float *rects, int total, long long *err1);
extern void g2d_gfx_gen_tex(void *data, const void *tex, int w, int h, int tex_unit, long long *err1);

#elif defined(G2D_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_H */