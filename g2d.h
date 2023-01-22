#ifndef G2D_H
#define G2D_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(G2D_WIN32)

typedef unsigned long g2d_ul_t;
typedef unsigned int g2d_ui_t;
extern void g2d_free(void *data);
extern void g2d_to_tstr(void **str, void *go_cstr, int *err_num);
extern void g2d_init(int *err_num, g2d_ul_t *err_win32);
extern void g2d_window_create(void **data, int go_obj, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void *t, int *err_num, g2d_ul_t *err_win32);
extern void g2d_window_props(void *data, int *x, int *y, int *w, int *h, int *wn, int *hn, int *wx, int *hx, int *b, int *d, int *r, int *f, int *l);
extern void g2d_process_messages();
extern void g2d_post_message(int *err_num, g2d_ul_t *err_win32);
extern void g2d_quit_message_queue();

#elif defined(G2D_LINUX)
#endif

#ifdef __cplusplus
}
#endif

#endif /* G2D_H */