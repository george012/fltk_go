#include "text.h"

#include <FL/Fl_Text_Display.H>

#include <FL/Fl_Text_Editor.H>

#include <FL/Fl_Text_Buffer.H>

#include "event_handler.h"


class GText_Display : public EventHandler<Fl_Text_Display> {
public:
  GText_Display(int x, int y, int w, int h, const char* label)
    : EventHandler<Fl_Text_Display>(x, y, w, h, label) {}
};

GText_Display *go_fltk_new_TextDisplay(int x, int y, int w, int h, const char *text) {
  return new GText_Display(x, y, w, h, text);
}

void go_fltk_TextDisplay_set_buffer(GText_Display *d, Fl_Text_Buffer *buf) {
  d->buffer(buf);
}

void go_fltk_TextDisplay_set_wrap_mode(GText_Display *b, int wrap, int wrapMargin) {
  b->wrap_mode(wrap, wrapMargin);
}

Fl_Text_Buffer *go_fltk_TextDisplay_buffer(GText_Display *d) {
  return d->buffer();
}

class GText_Editor : public EventHandler<Fl_Text_Editor> {
public:
  GText_Editor(int x, int y, int w, int h, const char* label)
    : EventHandler<Fl_Text_Editor>(x, y, w, h, label) {}
};

GText_Editor *go_fltk_new_TextEditor(int x, int y, int w, int h, const char *text) {
  return new GText_Editor(x, y, w, h, text);
}

Fl_Text_Buffer *go_fltk_new_TextBuffer(void) {
  return new Fl_Text_Buffer;
}

void go_fltk_TextBuffer_set_text(Fl_Text_Buffer *b, const char *txt) {
  b->text(txt);
}

void go_fltk_TextBuffer_append(Fl_Text_Buffer *b, const char *txt) {
  b->append(txt);
}

const char *go_fltk_TextBuffer_text(Fl_Text_Buffer *b) {
  return b->text();
}