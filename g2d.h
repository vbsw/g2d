#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WIN32)
typedef unsigned long g2d_ul_t;
typedef unsigned int g2d_ui_t;
typedef const char* g2d_lpcstr;
extern void g2d_error(void *err, int *err_num, g2d_ul_t *err_win32, char **err_str);
extern void g2d_error_free(void *err);
extern void *g2d_string_new(void **str, void *go_cstr);
extern void g2d_string_free(void *str);
extern void *g2d_init();
extern void *g2d_process_events();
extern void g2d_err_static_set(int go_obj);
extern void *g2d_window_create(void **data, int go_obj, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void *t);
extern void *g2d_window_show(void *data);
extern void *g2d_window_destroy(void *data, void **err);
extern void g2d_window_props(void *data, int *x, int *y, int *w, int *h, int *wn, int *hn, int *wx, int *hx, int *b, int *d, int *r, int *f, int *l);
extern void g2d_client_pos_set(void *data, int x, int y);
extern void g2d_client_size_set(void *data, int width, int height);
extern void *g2d_client_pos_apply(void *data);
extern void *g2d_client_move(void *data);
extern void g2d_window_style_set(void *data, int wn, int hn, int wx, int hx, int b, int d, int r, int l);
extern void g2d_client_restore_bak(void *data);
extern void *g2d_post_close(void *data);
extern void *g2d_post_update(void *data);
extern void *g2d_window_title_set(void *data, void *title);
extern void *g2d_mouse_pos_set(void *data, int x, int y);
extern void *g2d_window_fullscreen_set(void *data);
/*
extern void g2d_free(void *data);
extern void g2d_window_allocate(void **data, void **err);
extern void g2d_window_free(void *data, void **err);
extern void g2d_window_init_dummy(void *data, void **err);
extern void g2d_window_set_wgl_functions(void *data, void *cpf, void *cca);
extern void g2d_window_create(void *data, void **err);
extern void g2d_context_make_current(void *data, void **err);
extern void g2d_context_release(void *data, void **err);
extern void g2d_context_swap_buffers(void *data, void **err);
extern void g2d_window_context(void *data, void **ctx);
extern void g2d_window_show(void *data, void **err);
extern void g2d_window_props(void *data, int *x, int *y, int *w, int *h, int *wn, int *hn, int *wx, int *hx, int *b, int *d, int *r, int *f, int *l);
extern int g2d_window_funcs_avail(void *data);
extern int g2d_window_dt_func_avail(void *data);
*/
#elif defined(G2D_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_H */