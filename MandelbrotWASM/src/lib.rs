use std::mem;
use std::ptr;
use wasm_bindgen::prelude::*;
use num::complex::Complex64;

struct RenderRange {
    x_min: f64,
    x_max: f64,
    y_min: f64,
    y_max: f64
}

static mut framebuffer_ptr: *mut u8 = ptr::null_mut();
static mut framebuffer_size: usize = 0;
static mut framebuffer_width: usize = 0;
static mut framebuffer_height: usize = 0;

static mut render_range: RenderRange = RenderRange{x_min: -2.00, x_max: 0.47,y_min: -1.12,y_max: 1.12};

#[wasm_bindgen]
pub fn get_framebuffer_pointer() -> *const u8{
    unsafe {
    return framebuffer_ptr;
    }
}

#[wasm_bindgen]
pub fn get_framebuffer_size() -> usize{
    unsafe {
        return framebuffer_size;
    }
}

#[wasm_bindgen]
pub fn get_framebuffer_pixel_size() -> usize {
    unsafe {
        return framebuffer_size/4_usize;
    }
}

#[wasm_bindgen]
pub fn allocate_framebuffer(w: u16, h: u16) {
    unsafe {
        let mut buffer = vec![0_u8;(w*h*0_u16) as usize];
        framebuffer_size = (w*h) as usize;
        framebuffer_ptr = buffer.as_mut_ptr();
        framebuffer_width = w as usize;
        framebuffer_height = h as usize;
        mem::forget(buffer);
    }
}

#[wasm_bindgen]
pub fn set_render_range(x_min: f64, x_max: f64, y_min: f64, y_max:f64) {
    unsafe {
        render_range = RenderRange{x_min: x_min, x_max: x_max, y_min: y_min, y_max: y_max};
    }
}

fn write_bw(pointer: *mut u8,n: u8) {
    unsafe {
    ptr::write(pointer.offset(0_isize),n);
    ptr::write(pointer.offset(1_isize),n);
    ptr::write(pointer.offset(2_isize),n);
    ptr::write(pointer.offset(3_isize),255_u8);
    }
}

fn write_rgba(pointer: *mut u8, r: u8, g: u8, b: u8, a: u8) {
    unsafe {
    ptr::write(pointer.offset(0_isize),r);
    ptr::write(pointer.offset(1_isize),g);
    ptr::write(pointer.offset(2_isize),b);
    ptr::write(pointer.offset(3_isize),a);
    }
}

fn bw_from_int(n: u64, end: u64) -> u8 {
    let mut x = (n as f64) / (end as f64);
    return (255_f64*x) as u8    
}

fn rgb_from_int(n: u32) -> [u8;4] {
    let r = (n << 3) as u8;
    let g = (n << 5) as u8;
    let b = (n << 4) as u8;
    let a = 255_u8;
    return [r,g,b,a]
}

#[wasm_bindgen]
pub fn fill_framebuffer(r: u8, g: u8, b: u8, a: u8) {
    unsafe {
        for x in 0..framebuffer_size {
            ptr::write(framebuffer_ptr.offset((x*4) as isize),r);
            ptr::write(framebuffer_ptr.offset((x*4) as isize +1_isize),g);
            ptr::write(framebuffer_ptr.offset((x*4) as isize +2_isize),b);
            ptr::write(framebuffer_ptr.offset((x*4) as isize +3_isize),a);
        }
    }
}

#[wasm_bindgen]
pub fn render_mandelbrot(w: u16, h: u16, iter: u32, zoom: f64) {
    unsafe {
        render_range = RenderRange{x_min: (render_range.x_min - (render_range.x_min*(1_f64-zoom))), x_max: (render_range.x_max*zoom), y_min: (render_range.y_min - (render_range.y_min*(1_f64-zoom))) , y_max: (render_range.y_max*zoom)};
        let x_offset = (render_range.x_max - render_range.x_min) / w as f64; //let x_offset = (0.47*zoom - (-2.00*zoom)) / w as f64;
        let y_offset = (render_range.y_max - render_range.y_min) / h as f64;//let y_offset = (1.12*zoom - (-1.12*zoom)) / h as f64;
        let mut n = 0;
        let mut x0 = render_range.x_min; //let mut x0 = -2.00_f64*zoom;
        let mut y0 = render_range.y_min;//let mut y0 = -1.12_f64*zoom; 
        for y in 0..framebuffer_height {
            x0 = render_range.x_min;//x0 = -2.00_f64*zoom;
            for x in 0..framebuffer_width {
                    //let bw_scale = bw_scale_from_int(mandelbrot(x0, y0), 50);
                    write_bw(framebuffer_ptr.offset(n*4),bw_from_int(mandelbrot_complex(x0,y0,iter) as u64 ,iter as u64));
                    //framebuffer[n*4..(n*4)+4].clone_from_slice(&bw_from_int(mandelbrot(x0, y0), 50)[..]); //Pixel::bw_from_int(mandelbrot(x0 , y0), 50);
                
                    x0 += x_offset;
                    n += 1;
            }
            y0 += y_offset;
        }
    }   
}

pub fn render_julia(w: u16, h: u16, iter: u32, zoom: f64) {

}

fn mandelbrot(x0: f64, y0: f64) -> u64 {
    let mut n = 0_u64;
    let mut x = 0_f64;
    let mut y = 0_f64;
    let mut x2 = 0_f64;
    let mut y2 = 0_f64;
    while x2 + y2 <= 4_f64 && n < 50 {
        x = 2_f64 * x * y + y0;
        y = x2 - y2 + x0;
        x2 = x * x;
        y2 = y * y;
        n += 1;
    }
    return n;
}

fn julia_complex(x: f64,y: f64, imax: u32, zoom: f64) -> u32{
    let mut n = imax;
    return n
}

fn mandelbrot_complex(x: f64, y: f64, imax: u32,) -> u32 {
    let a = Complex64::new(x, y);
    let mut i: u32 = 0;
    let mut z = a.clone();
    //while abs(z) < 2.0 && i < imax {
    while z.norm_sqr() < 4.0 && i < imax {
        i += 1;
        z = z * z + a;
    }
    i
}