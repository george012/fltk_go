//go:build darwin

package fltk_go

// #cgo darwin,arm64,darwin,amd64 CXXFLAGS: -std=c++11
// #cgo darwin,arm64,darwin,amd64 CPPFLAGS: -I${SRCDIR}/include/darwin/universal -I${SRCDIR}/include/darwin/universal/FL/images -isysroot /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D_FILE_OFFSET_BITS=64 -D_THREAD_SAFE -D_REENTRANT
// #cgo darwin,arm64,darwin,amd64 LDFLAGS: ${SRCDIR}/lib/darwin/universal/libfltk_images.a ${SRCDIR}/lib/darwin/universal/libfltk_jpeg.a ${SRCDIR}/lib/darwin/universal/libfltk_png.a ${SRCDIR}/lib/darwin/universal/libfltk_z.a ${SRCDIR}/lib/darwin/universal/libfltk_gl.a -framework OpenGL ${SRCDIR}/lib/darwin/universal/libfltk_forms.a ${SRCDIR}/lib/darwin/universal/libfltk.a -lm -lpthread -framework Cocoa
import "C"
