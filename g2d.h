#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WIN32)

typedef unsigned long g2d_ul_t;
typedef unsigned int g2d_ui_t;
extern void g2d_free(void *data);
extern void g2d_to_tstr(void **str, void *go_cstr, size_t length, int *err_num);
extern void g2d_init(int *max_t_size, int *err_num, g2d_ul_t *err_win32);
extern void g2d_window_create(void **data, int cb_id, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void *t, int *err_num, g2d_ul_t *err_win32);
extern void g2d_window_show(void *data, int *err_num, g2d_ul_t *err_win32);
extern void g2d_window_destroy(void *data, int *err_num, g2d_ul_t *err_win32);
extern void g2d_window_props(void *data, int *x, int *y, int *w, int *h, int *wn, int *hn, int *wx, int *hx, int *b, int *d, int *r, int *f, int *l);
extern void g2d_process_messages();
extern void g2d_post_message(int *err_num, g2d_ul_t *err_win32);
extern void g2d_quit_message_queue();
extern void g2d_context_make_current(void *data, int *err_num, g2d_ul_t *err_win32);
extern void g2d_context_release(void *data, int *err_num, g2d_ul_t *err_win32);
extern void g2d_gfx_init(void *data, int *err_num, char **err_str);
extern void g2d_gfx_clear_bg(float r, float g, float b);
extern void g2d_gfx_swap_buffers(void *data, int *err_num, g2d_ul_t *err_win32);
extern void g2d_gfx_set_swap_interval(int interval);
extern void g2d_gfx_draw_rect(void *data, const char *enabled, const float *rects, int length, int active, int *err_num, char **err_str);
extern void g2d_gfx_draw_image(void *data, const char *enabled, const float *rects, int length, int active, int tex_id, int *err_num, char **err_str);
extern void g2d_gfx_set_view_size(void *data, int w, int h);
extern void g2d_gfx_gen_tex(void *data, const void *tex, int w, int h, int *tex_id, int *err_num);

#elif defined(G2D_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_H */