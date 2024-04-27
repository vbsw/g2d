#ifndef VBSW_G2D_WINDOW_H
#define VBSW_G2D_WINDOW_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct { void *set_func, *get_func; char *err_str; void **all; long long err1, err2; int list_len, list_cap, words_len, words_cap; } cdata_t;
extern void g2d_window_init(int pass, cdata_t *cdata);

extern void g2d_window_create(void **data, int cb_id, int x, int y, int w, int h, int wn, int hn, int wx, int hx, int b, int d, int r, int f, int l, int c, void *t, long long *err1, long long *err2);
extern void g2d_window_show(void *data, long long *err1, long long *err2);
extern void g2d_window_destroy(void *data, long long *err1, long long *err2);

extern void g2d_window_props(void *data, int *x, int *y, int *w, int *h, int *wn, int *hn, int *wx, int *hx, int *b, int *d, int *r, int *f, int *l);

extern void g2d_window_to_tstr(void **str, void *go_cstr, int *err_num);
extern void g2d_window_free(void *data);

#ifdef __cplusplus
}
#endif

#endif /* VBSW_G2D_WINDOW_H */