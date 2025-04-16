/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static GLuint shader_create(const GLenum shader_type, LPCSTR shader, const int err_a, const int err_b, long long *const err1, char **const err_nfo) {
	const GLuint prog_ref = glCreateShader(shader_type);
	if (prog_ref) {
		GLint compiled; glShaderSource(prog_ref, 1, &shader, NULL); glCompileShader(prog_ref);
		glGetShaderiv(prog_ref, GL_COMPILE_STATUS, &compiled);
		if (compiled == GL_FALSE) {
			GLsizei err_len; err1[0] = err_b; glGetShaderiv(prog_ref, GL_INFO_LOG_LENGTH, &err_len);
			if (err_len > 0) {
				err_nfo[0] = (char*)malloc(err_len);
				if (err_nfo[0])
					glGetShaderInfoLog(prog_ref, err_len, &err_len, (GLchar*)err_nfo[0]);
			}
			glDeleteShader(prog_ref);
		}
	} else {
		err1[0] = err_a;
	}
	return prog_ref;
}

static void shader_attach(const GLuint prog_ref, const GLuint shader_ref, const int err_a, const int err_b, long long *const err1) {
	glAttachShader(prog_ref, shader_ref);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_b;
	}
}

static void program_check(const GLuint prog_ref, const GLenum status, const int err, long long *const err1, char **const err_nfo) {
	GLint success; glGetProgramiv(prog_ref, status, &success);
	if (success == GL_FALSE) {
		GLsizei err_len; glGetProgramiv(prog_ref, GL_INFO_LOG_LENGTH, &err_len); err1[0] = err;
		if (err_len > 0) {
			err_nfo[0] = (char*)malloc(err_len);
			if (err_nfo[0])
				glGetProgramInfoLog(prog_ref, err_len, &err_len, err_nfo[0]);
		}
	}
}

static void prog_use(const GLuint id, const int err_a, const int err_b, long long *const err1) {
	glUseProgram(id);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_b;
	}
}

static GLuint rects_create(const GLuint vs_ref, const GLuint fs_ref, long long *const err1, char **const err_nfo) {
	if (err1[0] == 0) {
		const GLuint id = glCreateProgram();
		if (id) {
			shader_attach(id, vs_ref, G2D_ERR_1002006, G2D_ERR_1002007, err1);
			if (err1[0] == 0) {
				shader_attach(id, fs_ref, G2D_ERR_1002008, G2D_ERR_1002009, err1);
				if (err1[0] == 0) {
					glLinkProgram(id);
					program_check(id, GL_LINK_STATUS, G2D_ERR_1002010, err1, err_nfo);
				}
			}
		} else {
			err1[0] = G2D_ERR_1002011;
		}
		return id;
	}
	return 0;
}

static void bind_vao(const GLuint vao, const int err, long long *const err1) {
	if (err1[0] == 0) {
		glBindVertexArray(vao);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err;
		}
	}
}

static void bind_vbo(const GLuint vbo, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		glBindBuffer(GL_ARRAY_BUFFER, vbo);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
	}
}

static void bind_ebo(const GLuint ebo_ref, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, ebo_ref);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
	}
}

static void enable_attr(const GLint attr, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		glEnableVertexAttribArray(attr);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
	}
}

static GLint att_location(const GLuint prog_ref, LPCSTR const name, const int err, long long *const err1) {
	if (err1[0] == 0) {
		const GLint att_lc = glGetAttribLocation(prog_ref, name);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err;
		}
		return att_lc;
	}
	return -1;
}

static GLint unif_location(const GLuint prog_ref, LPCSTR const name, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		const GLint unif_lc = glGetUniformLocation(prog_ref, name);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
		return unif_lc;
	}
	return -1;
}

static void buffer_data(const GLenum target, const GLsizeiptr size, const void *const data, const GLenum usage, const int err_a, const int err_b, const int err_c, const int err_d, long long *const err1) {
	if (err1[0] == 0) {
		glBufferData(target, size, data, usage);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_c;
		} else if (err_enum == GL_OUT_OF_MEMORY) {
			err1[0] = err_d;
		}
	}
}

static void vertex_att_pointer(const GLuint index, const GLint size, const GLsizei stride, const void *const pointer, const int err_a, const int err_b, const int err_c, long long *const err1) {
	if (err1[0] == 0) {
		glVertexAttribPointer(index, size, GL_FLOAT, GL_FALSE, stride, pointer);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_b;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_c;
		}
	}
}

static void rects_enable(const GLuint prog_ref, const GLint unif_lc, const GLuint vao_ref, const GLuint vbo_ref, const GLfloat *const unif_data, long long *const err1) {
	prog_use(prog_ref, G2D_ERR_1002012, G2D_ERR_1002013, err1);
	if (err1[0] == 0) {
		bind_vao(vao_ref, G2D_ERR_1002014, err1);
		if (err1[0] == 0) {
			glUniform1fv(unif_lc, 16*3, unif_data);
			//glUniformMatrix4fv(unif_lc, 1, GL_FALSE, unif_data);
			bind_vbo(vbo_ref, G2D_ERR_1002015, G2D_ERR_1002016, err1);
		}
	}
}

static void buffer_sub_data(const GLsizeiptr size, const void *const data, const int err_a, const int err_b, const int err_c, long long *const err1) {
	glBufferSubData(GL_ARRAY_BUFFER, 0, size, data);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_b;
	} else if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_c;
	}
}

static void draw_elements(const GLsizei count, const int err_a, const int err_b, const int err_c, long long *const err1) {
	glDrawElements(GL_TRIANGLES, count, GL_UNSIGNED_INT, 0);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_b;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_c;
	}
}

static void bind_texture(const GLuint texture, const int err_a, const int err_b, const int err_c, long long *const err1) {
	glBindTexture(GL_TEXTURE_2D, texture);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_b;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_c;
	}
}

void g2d_gfx_gen_tex(void *const data, const void *const tex_data, const int w, const int h, const int tex_unit, long long *const err1) {
	GLuint texture; glGenTextures(1, &texture);
	glActiveTexture((GLenum)(GL_TEXTURE0+tex_unit));
	bind_texture(texture, G2D_ERR_1002052, G2D_ERR_1002053, G2D_ERR_1002054, err1);
	if (err1[0] == 0) {
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_BORDER);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_BORDER);
		//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
		//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
		//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
		//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
		glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, (GLsizei)w, (GLsizei)h, 0, GL_RGBA, GL_UNSIGNED_BYTE, tex_data);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = G2D_ERR_1002055;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = G2D_ERR_1002056;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = G2D_ERR_1002057;
		}
	}
}

void g2d_gfx_init(void *const data, long long *const err1, long long *const err2, char **const err_nfo) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wglMakeCurrent(wnd_data[0].wnd.dc, wnd_data[0].wnd.rc)) {
		const GLuint vs_id = shader_create(GL_VERTEX_SHADER, vs_rect_str, G2D_ERR_1002002, G2D_ERR_1002003, err1, err_nfo);
		if (err1[0] == 0) {
			const GLuint fs_id = shader_create(GL_FRAGMENT_SHADER, fs_rect_str, G2D_ERR_1002004, G2D_ERR_1002005, err1, err_nfo);
			if (err1[0] == 0) {
				const size_t length = 16000;
				wnd_data[0].rects.buf_max_len = (GLuint)length;
				wnd_data[0].rects.prog_ref = rects_create(vs_id, fs_id, err1, err_nfo);
				wnd_data[0].rects.att_lc[0] = att_location(wnd_data[0].rects.prog_ref, "in0", G2D_ERR_1002023, err1);
				wnd_data[0].rects.att_lc[1] = att_location(wnd_data[0].rects.prog_ref, "in1", G2D_ERR_1002023, err1);
				wnd_data[0].rects.att_lc[2] = att_location(wnd_data[0].rects.prog_ref, "in2", G2D_ERR_1002023, err1);
				wnd_data[0].rects.att_lc[3] = att_location(wnd_data[0].rects.prog_ref, "in3", G2D_ERR_1002023, err1);
				wnd_data[0].rects.unif_lc[0] = unif_location(wnd_data[0].rects.prog_ref, "tex00", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[1] = unif_location(wnd_data[0].rects.prog_ref, "tex01", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[2] = unif_location(wnd_data[0].rects.prog_ref, "tex02", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[3] = unif_location(wnd_data[0].rects.prog_ref, "tex03", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[4] = unif_location(wnd_data[0].rects.prog_ref, "tex04", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[5] = unif_location(wnd_data[0].rects.prog_ref, "tex05", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[6] = unif_location(wnd_data[0].rects.prog_ref, "tex06", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[7] = unif_location(wnd_data[0].rects.prog_ref, "tex07", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[8] = unif_location(wnd_data[0].rects.prog_ref, "tex08", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[9] = unif_location(wnd_data[0].rects.prog_ref, "tex09", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[10] = unif_location(wnd_data[0].rects.prog_ref, "tex10", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[11] = unif_location(wnd_data[0].rects.prog_ref, "tex11", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[12] = unif_location(wnd_data[0].rects.prog_ref, "tex12", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[13] = unif_location(wnd_data[0].rects.prog_ref, "tex13", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[14] = unif_location(wnd_data[0].rects.prog_ref, "tex14", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[15] = unif_location(wnd_data[0].rects.prog_ref, "tex15", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				wnd_data[0].rects.unif_lc[16] = unif_location(wnd_data[0].rects.prog_ref, "unif", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				if (err1[0] == 0) {
					GLuint objs[3]; glGenVertexArrays(1, objs); glGenBuffers(2, &objs[1]);
					wnd_data[0].rects.vao_ref = objs[0]; wnd_data[0].rects.vbo_ref = objs[1]; wnd_data[0].rects.ebo_ref = objs[2];
					bind_vao(wnd_data[0].rects.vao_ref, G2D_ERR_1002027, err1);
					enable_attr(wnd_data[0].rects.att_lc[0], G2D_ERR_1002028, G2D_ERR_1002029, err1);
					enable_attr(wnd_data[0].rects.att_lc[1], G2D_ERR_1002030, G2D_ERR_1002031, err1);
					enable_attr(wnd_data[0].rects.att_lc[2], G2D_ERR_1002032, G2D_ERR_1002033, err1);
					enable_attr(wnd_data[0].rects.att_lc[3], G2D_ERR_1002034, G2D_ERR_1002035, err1);
					bind_vbo(wnd_data[0].rects.vbo_ref, G2D_ERR_1002032, G2D_ERR_1002033, err1);
					buffer_data(GL_ARRAY_BUFFER, sizeof(GLfloat) * length * 4 * 16, NULL, GL_DYNAMIC_DRAW, G2D_ERR_1002034, G2D_ERR_1002035, G2D_ERR_1002036, G2D_ERR_1002037, err1);
					vertex_att_pointer(wnd_data[0].rects.att_lc[0], 4, sizeof(GLfloat) * 16, (void*)(sizeof(GLfloat) * 0), G2D_ERR_1002042, G2D_ERR_1002043, G2D_ERR_1002044, err1);
					vertex_att_pointer(wnd_data[0].rects.att_lc[1], 4, sizeof(GLfloat) * 16, (void*)(sizeof(GLfloat) * 4), G2D_ERR_1002042, G2D_ERR_1002043, G2D_ERR_1002044, err1);
					vertex_att_pointer(wnd_data[0].rects.att_lc[2], 4, sizeof(GLfloat) * 16, (void*)(sizeof(GLfloat) * 8), G2D_ERR_1002042, G2D_ERR_1002043, G2D_ERR_1002044, err1);
					vertex_att_pointer(wnd_data[0].rects.att_lc[3], 4, sizeof(GLfloat) * 16, (void*)(sizeof(GLfloat) * 12), G2D_ERR_1002042, G2D_ERR_1002043, G2D_ERR_1002044, err1);
					bind_ebo(wnd_data[0].rects.ebo_ref, G2D_ERR_1002048, G2D_ERR_1002049, err1);
					if (err1[0] == 0) {
						GLuint *indices = (GLuint*)malloc(sizeof(GLuint) * length * (3+3));
						if (indices) {
							wnd_data[0].rects.buffer = (GLfloat*)malloc(sizeof(GLfloat) * length * 4 * 16);
							if (wnd_data[0].rects.buffer) {
								size_t i;
								for (i = 0; i < length; i++) {
									const size_t offs = i * (3+3);
									const GLuint index = (GLuint) i * 4;
									indices[offs] = index; indices[offs+1] = index+1; indices[offs+2] = index+2; indices[offs+3] = index+2; indices[offs+4] = index+1; indices[offs+5] = index+3;
								}
								buffer_data(GL_ELEMENT_ARRAY_BUFFER, sizeof(GLuint) * length * (3+3), indices, GL_STATIC_DRAW, G2D_ERR_1002038, G2D_ERR_1002039, G2D_ERR_1002040, G2D_ERR_1002041, err1);
							} else {
								err1[0] = G2D_ERR_0000018;
							}
							free((void*)indices);
						} else {
							err1[0] = G2D_ERR_0000017;
						}
					}
				}
				glDeleteShader(fs_id);
			}
			glDeleteShader(vs_id);
		}
		if (err1[0] == 0) {
			glEnable(GL_BLEND);
			glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
		}
	} else {
		err1[0] = G2D_ERR_1002001, err2[0] = (long long)GetLastError();
	}
}

void g2d_gfx_release(void *const data, long long *const err1, long long *const err2) {
	if (!wglMakeCurrent(NULL, NULL))
		err1[0] = G2D_ERR_1002051, err2[0] = (long long)GetLastError();
}

void g2d_gfx_draw(void *const data, const int w, const int h, const int i, const float r, const float g, const float b,
	float **const buffs, const int *const bs, void **const procs, const int l, long long *const err1, long long *const err2) {
	int k;
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wnd_data[0].gfx.w != w || wnd_data[0].gfx.g != h) {
		wnd_data[0].gfx.w = w; wnd_data[0].gfx.h = h;
		wnd_data[0].gfx.unif_data[0] = 2.0f / (GLfloat)w;
		wnd_data[0].gfx.unif_data[5] = -2.0f / (GLfloat)h;
		glViewport((WORD)0, (WORD)0, (WORD)w, (WORD)h);
	}
	if (wnd_data[0].gfx.r != r || wnd_data[0].gfx.g != g || wnd_data[0].gfx.b != b) {
		wnd_data[0].gfx.r = r; wnd_data[0].gfx.g = g; wnd_data[0].gfx.b = b;
		glClearColor((GLfloat)r, (GLfloat)g, (GLfloat)b, 0.0);
	}
	if (wnd_data[0].gfx.i != i && wglSwapIntervalEXT) {
		wnd_data[0].gfx.i = i;
		wglSwapIntervalEXT(i);
	}
	glClear(GL_COLOR_BUFFER_BIT);
	for (k = 0; k < l && err1[0] == 0; k++) {
		gfx_draw_t *const draw = (gfx_draw_t*) procs[k];
		draw(data, buffs[k], bs[k], err1);
	}
	if (err1[0] == 0)
		if (!SwapBuffers(wnd_data[0].wnd.dc))
			err1[0] = G2D_ERR_1002050, err2[0] = (long long)GetLastError();
}

void g2d_gfx_draw_rectangles(void *const data, float *const rects, const int total, long long *const err1) {
	int rects_i, drawn;
	window_data_t *const wnd_data = (window_data_t*)data;
	const int length = (int)wnd_data[0].rects.buf_max_len;
	GLfloat *const buffer = wnd_data[0].rects.buffer;
	/* set dimensions (32=2*16) */
	for (rects_i = 16; rects_i < 48; rects_i++) {
		wnd_data[0].gfx.unif_data[rects_i] = rects[rects_i];
	}
	rects_enable(wnd_data[0].rects.prog_ref, wnd_data[0].rects.unif_lc[16], wnd_data[0].rects.vao_ref, wnd_data[0].rects.vbo_ref, wnd_data[0].gfx.unif_data, err1);
	/* set samplers (16) */
	for (rects_i = 0; rects_i < 16; rects_i++) {
		const int tex_unit = (int)rects[rects_i];
		if (tex_unit >= 0) {
			glUniform1i((GLint)wnd_data[0].rects.unif_lc[rects_i], (GLenum)tex_unit);
		//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
		//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
		}
	}
	for (rects_i = 0, drawn = 0; err1[0] == 0 && drawn < total; drawn += length) {
		int buf_i;
		const int limit = drawn + length > total ? total - drawn : length;
		for (buf_i = 0; buf_i < limit; rects_i++, buf_i++) {
			const int index = 48 + rects_i * 13; const int offs = buf_i * 4 * 16;
			const GLfloat x = rects[index], y = rects[index+1], w = rects[index+2], h = rects[index+3], r = rects[index+4], g = rects[index+5], b = rects[index+6], a = rects[index+7];
			const GLfloat sample = rects[index+8], tex_x = rects[index+9], tex_y = rects[index+10], tex_w = rects[index+11], tex_h = rects[index+12];
			buffer[offs+0] = x;
			buffer[offs+1] = y;
			buffer[offs+4] = r;
			buffer[offs+5] = g;
			buffer[offs+6] = b;
			buffer[offs+7] = a;
			buffer[offs+8] = sample;
			buffer[offs+12] = tex_x;
			buffer[offs+13] = tex_y;

			buffer[offs+16] = x + w;
			buffer[offs+17] = y;
			buffer[offs+20] = r;
			buffer[offs+21] = g;
			buffer[offs+22] = b;
			buffer[offs+23] = a;
			buffer[offs+24] = sample;
			buffer[offs+28] = tex_x + tex_w;
			buffer[offs+29] = tex_y;

			buffer[offs+32] = x;
			buffer[offs+33] = y + h;
			buffer[offs+36] = r;
			buffer[offs+37] = g;
			buffer[offs+38] = b;
			buffer[offs+39] = a;
			buffer[offs+40] = sample;
			buffer[offs+44] = tex_x;
			buffer[offs+45] = tex_y + tex_h;

			buffer[offs+48] = x + w;
			buffer[offs+49] = y + h;
			buffer[offs+52] = r;
			buffer[offs+53] = g;
			buffer[offs+54] = b;
			buffer[offs+55] = a;
			buffer[offs+56] = sample;
			buffer[offs+60] = tex_x + tex_w;
			buffer[offs+61] = tex_y + tex_h;
		}
		buffer_sub_data(sizeof(GLfloat) * limit * 4 * 16, buffer, G2D_ERR_1002020, G2D_ERR_1002021, G2D_ERR_1002022, err1);
		draw_elements(limit * 6, G2D_ERR_1002017, G2D_ERR_1002018, G2D_ERR_1002019, err1);
	}
}
