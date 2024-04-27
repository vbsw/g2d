#ifndef VBSW_G2D_MAINLOOP_H
#define VBSW_G2D_MAINLOOP_H

#ifdef __cplusplus
extern "C" {
#endif

extern void g2d_mainloop_process_messages();
extern void g2d_mainloop_post_custom(long long *err);
extern void g2d_mainloop_post_quit(long long *err);
extern void g2d_mainloop_clean_up();

#ifdef __cplusplus
}
#endif

#endif /* VBSW_G2D_MAINLOOP_H */