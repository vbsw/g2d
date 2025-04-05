/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static int keycode(const UINT message, const WPARAM wParam, const LPARAM lParam) {
	const int key = (int)(LOBYTE(HIWORD(lParam)));
	switch (key)
	{
	case 0:  return 0;
	case 1:  return 41;        // ESC         0x29
	case 2:  return 30;        // 1           0x1E
	case 3:  return 31;        // 2           0x1F
	case 4:  return 32;        // 3           0x20
	case 5:  return 33;        // 4           0x21
	case 6:  return 34;        // 5           0x22
	case 7:  return 35;        // 6           0x23
	case 8:  return 36;        // 7           0x24
	case 9:  return 37;        // 8           0x25
	case 10: return 38;        // 9           0x26
	case 11: return 39;        // 0           0x27
	case 12: return 45;        // -           0x2D
	case 13: return 46;        // =           0x2E
	case 14: return 42;        // DELETE      0x2A
	case 15: return 43;        // TAB         0x2B
	case 16: return 20;        // Q           0x14
	case 17: return 26;        // W           0x1A
	case 18: return 8;         // E           0x08
	case 19: return 21;        // R           0x15
	case 20: return 23;        // T           0x17
	case 21: return 28;        // Y           0x1C
	case 22: return 24;        // U           0x18
	case 23: return 12;        // I           0x0C
	case 24: return 18;        // O           0x12
	case 25: return 19;        // P           0x13
	case 26: return 47;        // [           0x2F
	case 27: return 48;        // ]           0x30
	case 28:
		if (HIBYTE(HIWORD(lParam)) & 0x1)
			return 88;         // pad ENTER   0x58
		return 40;             // board ENTER 0x28
	case 29:
		if (wParam == VK_CONTROL) {
			if (HIBYTE(HIWORD(lParam)) & 0x1)
				return 228;    // RCTRL       0xE4
			return 224;        // LCTRL       0xE0
		}
		return 0;
	case 30: return 4;         // A           0x04
	case 31: return 22;        // S           0x16
	case 32: return 7;         // D           0x07
	case 33: return 9;         // F           0x09
	case 34: return 10;        // G           0x0A
	case 35: return 11;        // H           0x0B
	case 36: return 13;        // J           0x0D
	case 37: return 14;        // K           0x0E
	case 38: return 15;        // L           0x0F
	case 39: return 51;        // ;           0x33
	case 40: return 52;        // '           0x34
	case 41: return 53;        // ^           0x35
	case 42: return 225;       // LSHIFT      0xE1
	case 43: return 50;        // ~           0x32
	case 44: return 29;        // Z           0x1D
	case 45: return 27;        // X           0x1B
	case 46: return 6;         // C           0x06
	case 47: return 25;        // V           0x19
	case 48: return 5;         // B           0x05
	case 49: return 17;        // N           0x11
	case 50: return 16;        // M           0x10
	case 51: return 54;        // ,           0x36
	case 52: return 55;        // .           0x37
	case 53:
		if (wParam == VK_DIVIDE)
			return 84;         // pad /       0x54
		return 56;             // /           0x38
	case 54: return 229;       // RSHIFT      0xE5
	case 55: return 85;        // pad *       0x55
	case 56:
		if (message == WM_SYSKEYDOWN || message == WM_SYSKEYUP)
			return 226;        // LALT        0xE2
		return 230;            // RALT        0xE6
	case 57: return 44;        // SPACE       0x2C
	case 58: return 57;        // CAPS        0x39
	case 59: return 58;        // F1          0x3A
	case 60: return 59;        // F2          0x3B
	case 61: return 60;        // F3          0x3C
	case 62: return 61;        // F4          0x3D
	case 63: return 62;        // F5          0x3E
	case 64: return 63;        // F6          0x3F
	case 65: return 64;        // F7          0x40
	case 66: return 65;        // F8          0x41
	case 67: return 66;        // F9          0x42
	case 68: return 67;        // F10         0x43
	case 69:
		if (wParam == VK_PAUSE)
			return 72;         // PAUSE       0x48
		return 83;             // pad LOCK    0x53
	case 70: return 71;        // SCROLL      0x47
	case 71:
		if (wParam == VK_HOME)
			return 74;         // HOME        0x4A
		return 95;             // pad 7       0x5F
	case 72:
		if (wParam == VK_UP)
			return 82;         // UP          0x52
		return 96;             // pad 8       0x60
	case 73:
		if (wParam == VK_PRIOR)
			return 75;         // PAGEUP      0x4B
		return 97;             // pad 9       0x61
	case 74: return 86;        // pad -       0x56
	case 75:
		if (wParam == VK_LEFT)
			return 80;         // LEFT        0x50
		return 92;             // pad 4       0x5C
	case 76: return 93;        // pad 5       0x5D
	case 77:
		if (wParam == VK_RIGHT)
			return 79;         // RIGHT       0x4F
		return 94;             // pad 6       0x5E
	case 78: return 87;        // pad +       0x57
	case 79:
		if (wParam == VK_END)
			return 77;         // END         0x4D
		return 89;             // pad 1       0x59
	case 80:
		if (wParam == VK_DOWN)
			return 81;         // DOWN        0x51
		return 90;             // pad 2       0x5A
	case 81:
		if (wParam == VK_NEXT)
			return 78;         // PAGEDOWN    0x4E
		return 91;             // pad 3       0x5B
	case 82: return 73;        // INSERT      0x49
	case 83:
		if (wParam == VK_DELETE)
			return 76;         // DELETE F    0x4C
		return 99;             // pad DELETE  0x63
	case 84: return 0;
	case 85: return 0;
	case 86: return 100;       // |           0x64
	case 87: return 68;        // F11         0x44
	case 88: return 69;        // F12         0x45
	case 89: return 0;         // LWIN        0xE3
	case 90: return 0;         // RWIN        0xE7
	case 91: return 0;
	case 92: return 0;
	case 93: return 118;       // MENU        0x76
	}
	return key;
}

static BOOL key_down_process(window_data_t *const wnd_data, const UINT message, const WPARAM wParam, const LPARAM lParam) {
	const int code = keycode(message, wParam, lParam);
	if (code) {
		g2dKeyDown(wnd_data[0].cb_id, code, wnd_data[0].key_repeated[code]++);
		return TRUE;
	}
	return FALSE;
}

static BOOL key_up_process(window_data_t *const wnd_data, const UINT message, const WPARAM wParam, const LPARAM lParam) {
	const int code = keycode(message, wParam, lParam);
	if (code) {
		wnd_data[0].key_repeated[code] = 0;
		g2dKeyUp(wnd_data[0].cb_id, code);
		return TRUE;
	}
	return FALSE;
}
