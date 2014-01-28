// Copyright 2013 Benoît Amiaux. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rez

type Config struct {
	depth      int
	input      int
	output     int
	vertical   bool
	interlaced bool
}

type Resizer interface {
	Resize(dst, src []byte, dpitch, spitch int, width, height int)
}

type Scaler func(taps, width, height int, coeffs []int16, offset []int,
	dst, src []byte, dpitch, spitch int)

type Context struct {
	cfg     Config
	kernels []Kernel
	scaler  Scaler
}

func NewResize(cfg *Config, filter Filter) Resizer {
	ctx := Context{
		cfg:    *cfg,
		scaler: h8scaleN,
	}
	ctx.cfg.depth = 8 // only 8-bit for now
	ctx.kernels = []Kernel{makeKernel(&ctx.cfg, filter, 0)}
	if cfg.vertical {
		ctx.scaler = v8scaleN
		if cfg.interlaced {
			ctx.kernels = append(ctx.kernels, makeKernel(&ctx.cfg, filter, 1))
		}
	}
	return &ctx
}

func (c *Context) Resize(dst, src []byte, dpitch, spitch int, width, height int) {
	field := bin(c.cfg.vertical && c.cfg.interlaced)
	dwidth := c.cfg.output
	dheight := height
	if c.cfg.vertical {
		dwidth = width
		dheight = c.cfg.output >> field
	}
	for i, k := range c.kernels[:1+field] {
		c.scaler(k.size, dwidth, dheight, k.coeffs, k.offsets,
			dst[dpitch*i:], src[spitch*i:], dpitch<<field, spitch<<field)
	}
}