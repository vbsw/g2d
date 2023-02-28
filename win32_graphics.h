/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

/* orthographic projection */
static const float projection_mat[4*4] = { 2.0f / 1.0f, 0.0f, 0.0f, 0.0f, 0.0f, -2.0f / 1.0f, 0.0f, 0.0f, 0.0f, 0.0f, -1.0f, 0.0f, -1.0f, 1.0f, 0.0f, 1.0f };
static LPCSTR const vs_rect_str = "#version 130\nin vec2 positionIn; in vec4 colorIn; out vec4 fragementColor; uniform mat4 projection = mat4(1.0); void main() { gl_Position = projection * vec4(positionIn, 1.0, 1.0); fragementColor = colorIn; }";
static LPCSTR const fs_rect_str = "#version 130\nin vec4 fragementColor; out vec4 color; void main() { color = fragementColor; }";


static LPCSTR const vs_image_str = "#version 130\n\
in vec2 positionIn; \
in vec4 colorIn; \
in vec2 texCoordIn; \
out vec4 fragementColor; \
out vec2 fragementTexCoord; \
uniform mat4 projection = mat4(1.0); \
void main() { \
gl_Position = projection * vec4(positionIn, 1.0, 1.0); \
fragementColor = colorIn; \
fragementTexCoord = texCoordIn; }";

static LPCSTR const fs_image_str = "#version 130\n\
in vec4 fragementColor; \
in vec2 fragementTexCoord; \
out vec4 color; \
uniform sampler2D textureIn;\
void main() { \
color = fragementColor; \
color = texture(textureIn, fragementTexCoord); \
}";


static GLuint shader_create(const GLenum shader_type, LPCSTR shader, const int err_a, const int err_b, int *const err_num, char **const err_str) {
	const GLuint id = glCreateShader(shader_type);
	if (id) {
		GLint compiled; glShaderSource(id, 1, &shader, NULL); glCompileShader(id);
		glGetShaderiv(id, GL_COMPILE_STATUS, &compiled);
		if (compiled == GL_FALSE) {
			GLsizei err_len; err_num[0] = err_b; glGetShaderiv(id, GL_INFO_LOG_LENGTH, &err_len);
			if (err_len > 0) {
				err_str[0] = (char*)malloc(err_len);
				if (err_str[0])
					glGetShaderInfoLog(id, err_len, &err_len, (GLchar*)err_str[0]);
			}
			glDeleteShader(id);
		}
	} else {
		err_num[0] = err_a;
	}
	return id;
}

static void shader_attach(const GLuint prog_id, const GLuint shader_id, const int err_a, const int err_b, int *const err_num) {
	glAttachShader(prog_id, shader_id);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err_num[0] = err_b;
	}
}

static void program_check(const GLuint prog_id, const GLenum status, const int err, int *const err_num, char **const err_str) {
	GLint success; glGetProgramiv(prog_id, status, &success);
	if (success == GL_FALSE) {
		GLsizei err_len; glGetProgramiv(prog_id, GL_INFO_LOG_LENGTH, &err_len); err_num[0] = err;
		if (err_len > 0) {
			err_str[0] = (char*)malloc(err_len);
			if (err_str[0])
				glGetProgramInfoLog(prog_id, err_len, &err_len, err_str[0]);
		}
	}
}

static void prog_use(const GLuint id, const int err_a, const int err_b, int *const err_num) {
	glUseProgram(id);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err_num[0] = err_b;
	}
}

static GLuint rect_prog_create(const GLuint vs_id, const GLuint fs_id, int *const err_num, char **const err_str) {
	if (err_num[0] == 0) {
		const GLuint id = glCreateProgram();
		if (id) {
			shader_attach(id, vs_id, 1005, 1006, err_num);
			if (err_num[0] == 0) {
				shader_attach(id, fs_id, 1007, 1008, err_num);
				if (err_num[0] == 0) {
					glLinkProgram(id);
					program_check(id, GL_LINK_STATUS, 1009, err_num, err_str);
				}
			}
		} else {
			err_num[0] = 1004;
		}
		return id;
	}
	return 0;
}

static GLuint image_prog_create(const GLuint vs_id, const GLuint fs_id, int *const err_num, char **const err_str) {
	if (err_num[0] == 0) {
		const GLuint id = glCreateProgram();
		if (id) {
			shader_attach(id, vs_id, 1105, 1106, err_num);
			if (err_num[0] == 0) {
				shader_attach(id, fs_id, 1107, 1108, err_num);
				if (err_num[0] == 0) {
					glLinkProgram(id);
					program_check(id, GL_LINK_STATUS, 1109, err_num, err_str);
				}
			}
		} else {
			err_num[0] = 1104;
		}
		return id;
	}
	return 0;
}

static void enable_attr(const GLint attr, const int err_a, const int err_b, int *const err_num) {
	if (err_num[0] == 0) {
		glEnableVertexAttribArray(attr);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = err_b;
		}
	}
}

static GLint att_location(const GLuint prog_id, LPCSTR const name, const int err, int *const err_num) {
	if (err_num[0] == 0) {
		const GLint att_id = glGetAttribLocation(prog_id, name);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = err;
		}
		return att_id;
	}
	return -1;
}

static GLint unf_location(const GLuint prog_id, LPCSTR const name, const int err_a, const int err_b, int *const err_num) {
	if (err_num[0] == 0) {
		const GLint att_id = glGetUniformLocation(prog_id, name);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = err_b;
		}
		return att_id;
	}
	return -1;
}

static void bind_vao(const GLuint vao, const int err, int *const err_num) {
	if (err_num[0] == 0) {
		glBindVertexArray(vao);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = err;
		}
	}
}

static void bind_vbo(const GLuint vbo, const int err_a, const int err_b, int *const err_num) {
	if (err_num[0] == 0) {
		glBindBuffer(GL_ARRAY_BUFFER, vbo);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err_num[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = err_b;
		}
	}
}

static void bind_ebo(const GLuint ebo, const int err_a, const int err_b, int *const err_num) {
	if (err_num[0] == 0) {
		glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, ebo);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err_num[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = err_b;
		}
	}
}

static void bind_texture(const GLuint texture, int *const err_num) {
	glBindTexture(GL_TEXTURE_2D, texture);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err_num[0] = 1170;
	} else if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = 1171;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err_num[0] = 1172;
	}
}

static void buffer_data(const GLenum target, const GLsizeiptr size, const void *const data, const GLenum usage, const int err_a, const int err_b, const int err_c, const int err_d, int *const err_num) {
	if (err_num[0] == 0) {
		glBufferData(target, size, data, usage);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err_num[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = err_b;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = err_c;
		} else if (err_enum == GL_OUT_OF_MEMORY) {
			err_num[0] = err_d;
		}
	}
}

static void vertex_att_pointer(const GLuint index, const GLint size, const GLsizei stride, const void *const pointer, const int err_a, const int err_b, const int err_c, int *const err_num) {
	if (err_num[0] == 0) {
		glVertexAttribPointer(index, size, GL_FLOAT, GL_FALSE, stride, pointer);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = err_a;
		} else if (err_enum == GL_INVALID_ENUM) {
			err_num[0] = err_b;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = err_c;
		}
	}
}

static void rect_prog_enable(const GLuint prog_id, const GLint proj_unif, const GLuint vao, const GLuint vbo, const float *const projection_mat, int *const err_num) {
	prog_use(prog_id, 1200, 1201, err_num);
	if (err_num[0] == 0) {
		bind_vao(vao, 1202, err_num);
		if (err_num[0] == 0) {
			glUniformMatrix4fv(proj_unif, 1, GL_FALSE, projection_mat);
			bind_vbo(vbo, 1203, 1204, err_num);
		}
	}
}

static void buffer_sub_data(const GLsizeiptr size, const void *const data, const int err_a, const int err_b, const int err_c, int *const err_num) {
	glBufferSubData(GL_ARRAY_BUFFER, 0, size, data);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err_num[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err_num[0] = err_b;
	} else if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = err_c;
	}
}

static void draw_elements(const GLsizei count, const int err_a, const int err_b, const int err_c, int *const err_num) {
	glDrawElements(GL_TRIANGLES, count, GL_UNSIGNED_INT, 0);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err_num[0] = err_a;
	} else if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = err_b;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err_num[0] = err_c;
	}
}

static void rect_init(window_data_t *const wnd_data, int *const err_num, char **const err_str) {
	if (err_num[0] == 0) {
		const GLuint vs_id = shader_create(GL_VERTEX_SHADER, vs_rect_str, 1000, 1001, err_num, err_str);
		if (err_num[0] == 0) {
			const GLuint fs_id = shader_create(GL_FRAGMENT_SHADER, fs_rect_str, 1002, 1003, err_num, err_str);
			if (err_num[0] == 0) {
				const size_t size = 16000;
				wnd_data[0].rect_prog.max_size = (GLuint)size;
				wnd_data[0].rect_prog.id = rect_prog_create(vs_id, fs_id, err_num, err_str);
				wnd_data[0].rect_prog.pos_att = att_location(wnd_data[0].rect_prog.id, "positionIn", 1010, err_num);
				wnd_data[0].rect_prog.col_att = att_location(wnd_data[0].rect_prog.id, "colorIn", 1011, err_num);
				wnd_data[0].rect_prog.proj_unif = unf_location(wnd_data[0].rect_prog.id, "projection", 1012, 1013, err_num);
				bind_vao(wnd_data[0].rect_prog.vao, 1014, err_num);
				enable_attr(wnd_data[0].rect_prog.pos_att, 1015, 1016, err_num);
				enable_attr(wnd_data[0].rect_prog.col_att, 1017, 1018, err_num);
				bind_vbo(wnd_data[0].rect_prog.vbo, 1019, 1020, err_num);
				buffer_data(GL_ARRAY_BUFFER, sizeof(float) * size * 4 * (2+4), NULL, GL_DYNAMIC_DRAW, 1021, 1022, 1023, 1024, err_num);
				vertex_att_pointer(wnd_data[0].rect_prog.pos_att, 2, sizeof(float) * (2+4), (void*)(sizeof(float) * 0), 1025, 1026, 1027, err_num);
				vertex_att_pointer(wnd_data[0].rect_prog.col_att, 4, sizeof(float) * (2+4), (void*)(sizeof(float) * 2), 1028, 1029, 1030, err_num);
				bind_ebo(wnd_data[0].rect_prog.ebo, 1031, 1032, err_num);
				if (err_num[0] == 0) {
					unsigned int *indices = (unsigned int*)malloc(sizeof(unsigned int) * size * (2+4));
					if (indices) {
						wnd_data[0].rect_prog.buffer = (float*)malloc(sizeof(float) * size * 4 * (2+4));
						if (wnd_data[0].rect_prog.buffer) {
							size_t i;
							for (i = 0; i < size; i++) {
								const size_t offs = i * (3+3);
								const size_t index = i * 4;
								indices[offs] = index; indices[offs+1] = index+1; indices[offs+2] = index+2; indices[offs+3] = index+2; indices[offs+4] = index+1; indices[offs+5] = index+3;
							}
							buffer_data(GL_ELEMENT_ARRAY_BUFFER, sizeof(unsigned int) * size * (3+3), indices, GL_STATIC_DRAW, 1035, 1036, 1037, 1038, err_num);
						} else {
							err_num[0] = 1034;
						}
						free((void*)indices);
					} else {
						err_num[0] = 1033;
					}
				}
				glDeleteShader(fs_id);
			}
			glDeleteShader(vs_id);
		}
	}
}

static void image_init(window_data_t *const wnd_data, int *const err_num, char **const err_str) {
	if (err_num[0] == 0) {
		const GLuint vs_id = shader_create(GL_VERTEX_SHADER, vs_image_str, 1100, 1101, err_num, err_str);
		if (err_num[0] == 0) {
			const GLuint fs_id = shader_create(GL_FRAGMENT_SHADER, fs_image_str, 1102, 1103, err_num, err_str);
			if (err_num[0] == 0) {
				const size_t size = 16000;
				wnd_data[0].image_prog.max_size = (GLuint)size;
				wnd_data[0].image_prog.id = image_prog_create(vs_id, fs_id, err_num, err_str);
				wnd_data[0].image_prog.pos_att = att_location(wnd_data[0].image_prog.id, "positionIn", 1110, err_num);
				wnd_data[0].image_prog.col_att = att_location(wnd_data[0].image_prog.id, "colorIn", 1111, err_num);
				wnd_data[0].image_prog.proj_unif = unf_location(wnd_data[0].image_prog.id, "projection", 1112, 1113, err_num);
				bind_vao(wnd_data[0].image_prog.vao, 1114, err_num);
				enable_attr(wnd_data[0].image_prog.pos_att, 1115, 1116, err_num);
				enable_attr(wnd_data[0].image_prog.col_att, 1117, 1118, err_num);
				bind_vbo(wnd_data[0].image_prog.vbo, 1119, 1120, err_num);
				buffer_data(GL_ARRAY_BUFFER, sizeof(float) * size * 4 * (2+4+2), NULL, GL_DYNAMIC_DRAW, 1121, 1122, 1123, 1124, err_num);
				vertex_att_pointer(wnd_data[0].image_prog.pos_att, 2, sizeof(float) * (2+4+2), (void*)(sizeof(float) * 0), 1125, 1126, 1127, err_num);
				vertex_att_pointer(wnd_data[0].image_prog.col_att, 4, sizeof(float) * (2+4+2), (void*)(sizeof(float) * 2), 1128, 1129, 1130, err_num);
				vertex_att_pointer(wnd_data[0].image_prog.tex_att, 2, sizeof(float) * (2+4+2), (void*)(sizeof(float) * (2+4)), 1031, 1032, 1033, err_num);
				bind_ebo(wnd_data[0].image_prog.ebo, 1134, 1135, err_num);
				if (err_num[0] == 0) {
					unsigned int *indices = (unsigned int*)malloc(sizeof(unsigned int) * size * (2+4+2));
					if (indices) {
						wnd_data[0].image_prog.buffer = (float*)malloc(sizeof(float) * size * 4 * (2+4+2));
						if (wnd_data[0].image_prog.buffer) {
							size_t i;
							for (i = 0; i < size; i++) {
								const size_t offs = i * (3+3);
								const size_t index = i * 4;
								indices[offs] = index; indices[offs+1] = index+1; indices[offs+2] = index+2; indices[offs+3] = index+2; indices[offs+4] = index+1; indices[offs+5] = index+3;
							}
							buffer_data(GL_ELEMENT_ARRAY_BUFFER, sizeof(unsigned int) * size * (3+3), indices, GL_STATIC_DRAW, 1138, 1139, 1140, 1141, err_num);
						} else {
							err_num[0] = 1137;
						}
						free((void*)indices);
					} else {
						err_num[0] = 1136;
					}
				}
				glDeleteShader(fs_id);
			}
			glDeleteShader(vs_id);
		}
	}
}

void g2d_gfx_init(void *const data, int *const err_num, char **const err_str) {
	GLuint objs[6]; window_data_t *const wnd_data = (window_data_t*)data;
	glGenVertexArrays(2, objs); glGenBuffers(4, &objs[2]);
	wnd_data[0].rect_prog.vao = objs[0];
	wnd_data[0].rect_prog.vbo = objs[2];
	wnd_data[0].rect_prog.ebo = objs[3];
	wnd_data[0].image_prog.vao = objs[1];
	wnd_data[0].image_prog.vbo = objs[4];
	wnd_data[0].image_prog.ebo = objs[5];
	rect_init(wnd_data, err_num, err_str);
	if (err_num[0] == 0) {
		image_init(wnd_data, err_num, err_str);
		if (err_num[0] == 0) {
			glEnable(GL_BLEND);
			glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
			memcpy(&wnd_data[0].projection_mat, projection_mat, sizeof(float)*(4*4));
		}
	}
}

void g2d_gfx_clear_bg(const float r, const float g, const float b) {
	glClearColor((GLclampf)r, (GLclampf)g, (GLclampf)b, 0.0);
	glClear(GL_COLOR_BUFFER_BIT);
}

void g2d_gfx_swap_buffers(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	//glFlush();
	//glFinish();
	if (!SwapBuffers(((window_data_t*)data)[0].wnd.ctx.dc)) {
		err_num[0] = 61; err_win32[0] = (g2d_ul_t)GetLastError();
	}
}

void g2d_gfx_set_swap_interval(const int interval) {
	wglSwapIntervalEXT(interval);
}

void g2d_gfx_draw_rect(void *const data, const char *const enabled, const float *const rects, const int length, const int active, int *const err_num, char **const err_str) {
	if (active > 0) {
		int i, drawn;
		window_data_t *const wnd_data = (window_data_t*)data;
		const int size = (int)wnd_data[0].rect_prog.max_size;
		float *const buffer = wnd_data[0].rect_prog.buffer;
		rect_prog_enable(wnd_data[0].rect_prog.id, wnd_data[0].rect_prog.proj_unif, wnd_data[0].rect_prog.vao, wnd_data[0].rect_prog.vbo, wnd_data[0].projection_mat, err_num);
		for (i = 0, drawn = 0; err_num[0] == 0 && drawn < active; drawn += size) {
			int k;
			const int limit = drawn + size > active ? active - drawn : size;
			for (k = 0; k < limit; i++) {
				if (enabled[i]) {
					const int offs = k * 4 * (2+4); const int index = i * 20;
					const float x = rects[index], y = rects[index+1], w = rects[index+2], h = rects[index+3];
					buffer[offs] = x;
					buffer[offs+1] = y;
					buffer[offs+2] = rects[index+4]; // r
					buffer[offs+3] = rects[index+5]; // g
					buffer[offs+4] = rects[index+6]; // b
					buffer[offs+5] = rects[index+7]; // a
					buffer[offs+6] = x + w;
					buffer[offs+7] = y;
					buffer[offs+8] = rects[index+8];
					buffer[offs+9] = rects[index+9];
					buffer[offs+10] = rects[index+10];
					buffer[offs+11] = rects[index+11];
					buffer[offs+12] = x;
					buffer[offs+13] = y + h;
					buffer[offs+14] = rects[index+12];
					buffer[offs+15] = rects[index+13];
					buffer[offs+16] = rects[index+14];
					buffer[offs+17] = rects[index+15];
					buffer[offs+18] = x + w;
					buffer[offs+19] = y + h;
					buffer[offs+20] = rects[index+16];
					buffer[offs+21] = rects[index+17];
					buffer[offs+22] = rects[index+18];
					buffer[offs+23] = rects[index+19];
					k++;
				}
			}
			buffer_sub_data(sizeof(float) * limit * 4 * (2+4), buffer, 1205, 1206, 1207, err_num);
			draw_elements(limit * 6, 1208, 1209, 1210, err_num);
		}
	}
}

void g2d_gfx_draw_image(void *const data, const char *const enabled, const float *const images, const int length, const int active, int *const err_num, char **const err_str) {
	if (active > 0) {
		int i, drawn;
		window_data_t *const wnd_data = (window_data_t*)data;
		const int size = (int)wnd_data[0].rect_prog.max_size;
		float *const buffer = wnd_data[0].rect_prog.buffer;
		rect_prog_enable(wnd_data[0].rect_prog.id, wnd_data[0].rect_prog.proj_unif, wnd_data[0].rect_prog.vao, wnd_data[0].rect_prog.vbo, wnd_data[0].projection_mat, err_num);
		for (i = 0, drawn = 0; err_num[0] == 0 && drawn < active; drawn += size) {
			int k;
			const int limit = drawn + size > active ? active - drawn : size;
			for (k = 0; k < limit; i++) {
				if (enabled[i]) {
					const int offs = k * 4 * (2+4+2); const int index = i * 24;
					const float x = images[index], y = images[index+1], w = images[index+2], h = images[index+3];
					const float tx = images[index+4], ty = images[index+5], tw = images[index+6], th = images[index+7];
					buffer[offs] = x;
					buffer[offs+1] = y;
					buffer[offs+2] = images[index+8]; // r
					buffer[offs+3] = images[index+9]; // g
					buffer[offs+4] = images[index+10]; // b
					buffer[offs+5] = images[index+11]; // a
					buffer[offs+6] = tx;
					buffer[offs+7] = ty;
					buffer[offs+8] = x + w;
					buffer[offs+9] = y;
					buffer[offs+10] = images[index+12];
					buffer[offs+11] = images[index+13];
					buffer[offs+12] = images[index+14];
					buffer[offs+13] = images[index+15];
					buffer[offs+14] = tx + tw;
					buffer[offs+15] = ty;
					buffer[offs+16] = x;
					buffer[offs+17] = y + h;
					buffer[offs+18] = images[index+16];
					buffer[offs+19] = images[index+17];
					buffer[offs+20] = images[index+18];
					buffer[offs+21] = images[index+19];
					buffer[offs+22] = tx;
					buffer[offs+23] = ty + th;
					buffer[offs+24] = x + w;
					buffer[offs+25] = y + h;
					buffer[offs+26] = images[index+20];
					buffer[offs+27] = images[index+21];
					buffer[offs+28] = images[index+22];
					buffer[offs+29] = images[index+23];
					buffer[offs+30] = tx + tw;
					buffer[offs+31] = ty + th;
					k++;
				}
			}
			buffer_sub_data(sizeof(float) * limit * 4 * (2+4+2), buffer, 1205, 1206, 1207, err_num);
			draw_elements(limit * 6, 1208, 1209, 1210, err_num);
		}
	}
}

void g2d_gfx_set_view_size(void *const data, const int w, const int h) {
	window_data_t *const wnd_data = (window_data_t*)data;
	glViewport((WORD)0, (WORD)0, (WORD)w, (WORD)h);
	wnd_data[0].projection_mat[0] = 2.0f / (float)w;
	wnd_data[0].projection_mat[5] = -2.0f / (float)h;
}

void g2d_gfx_gen_tex(void *const data, const void *const tex, const int w, const int h, int *const tex_id, int *const err_num) {
	GLuint textures;
	glGenTextures(1, &textures);
	bind_texture(textures, err_num);
	if (err_num[0] == 0) {
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_BORDER);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_BORDER);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
		glTexImage2D(GL_TEXTURE_2D, 0, 4, (GLsizei)w, (GLsizei)h, 0, GL_RGBA, GL_UNSIGNED_BYTE, tex);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_NO_ERROR) {
			tex_id[0] = (int)textures;
		} else if (err_enum == GL_INVALID_ENUM) {
			err_num[0] = 1173;
		} else if (err_enum == GL_INVALID_VALUE) {
			err_num[0] = 1174;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err_num[0] = 1175;
		}
	}
}
