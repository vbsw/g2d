#ifndef VBSW_G2D_NLTX_H
#define VBSW_G2D_NLTX_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct { void *set_func, *get_func; char *err_str; void **all; long long err1, err2; int list_len, list_cap, words_len, words_cap; } cdata_t;
extern void vbsw_nltx_init(int pass, cdata_t *cdata);
extern void vbsw_nltx_result_and_free(void *data, int *mts, int *mtu);

#ifdef __cplusplus
}
#endif

#endif /* VBSW_G2D_NLTX_H */