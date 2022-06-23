#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WIN32)
typedef unsigned long g2d_ul_t;
extern void g2d_error(void *err, int *err_num, g2d_ul_t *err_win32, char **err_str);
extern void g2d_error_free(void *err);
extern void g2d_init(void **err);
extern void g2d_process_events(void **err);
/*
extern void g2d_free(void *data);
extern void g2d_window_allocate(void **data, void **err);
extern void g2d_window_free(void *data, void **err);
extern void g2d_window_init_dummy(void *data, void **err);
extern void g2d_window_init_opengl30(void *data, int go_obj, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void **err);
extern void g2d_window_set_wgl_functions(void *data, void *cpf, void *cca);
extern void g2d_window_create(void *data, void **err);
extern void g2d_window_destroy(void *data, void **err);
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