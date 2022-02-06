# WasmOnCanvas

## About
This library aims to provide JavaScript bindinds to render patters on canvas with WebAssembly.

## Prerequisites

`wasm-pack` (https://rustwasm.github.io/wasm-pack/)

## Build

`wasm-pack build --target web`

## Releases

Pre compiled binaries can be found in the releases section.

## Example

```
<script type="module">
import init, {allocate_framebuffer,set_render_range,get_framebuffer_size,get_framebuffer_pointer,get_framebuffer_pixel_size,fill_framebuffer,render_mandelbrot} from './canvas_wasm.js';

const run = async () => {
  const wasm = await init("./canvas_wasm_bg.wasm"); // Loads WASM file
  
  const width = 200;
  const height = 200;

  const canvas = document.querySelector("canvas");
  const canvas_context = canvas.getContext("2d");
  const canvas_image = canvas_context.createImageData(width,height)

  allocate_framebuffer(width,height);  // Allocates the framebuffer

  const wasm_memory = new Uint8Array(wasm.memory.buffer); // Gets Uint8 Array with memory, must be done after framebuffer allocation

  const framebuffer_ptr = get_framebuffer_pointer(); // Gets the pointer for the framebuffer memory

  const image_data = canvas_context.createImageData(width,height); // Creates image data

  canvas_context.clearRect(0,0,width,height); // Clears the canvas

  render_mandelbrot(width,height,50,1); // Renders the mandelbrot set to the framebuffer

  image_array = wasm_memory.slice(framebuffer_ptr,framebuffer_ptr+(width*height*4)); // Gets the framebuffer array
  image_data.data.set(image_array); // Creates image data from byte array
  canvas_context.putImageData(image_data,0,0) // Puts the image data in the canvas framebuffer.

}

run();

</script>
```