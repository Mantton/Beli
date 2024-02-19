import type { MetaFunction } from "@remix-run/node";
import { useEffect, useRef } from "react";

export const meta: MetaFunction = () => {
  return [
    { title: "Beli - Colors" },
    { name: "description", content: "Beli - An Interactive Canvas" },
  ];
};

export default function Board() {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const CELL_SIZE = 75; // Cell Size

  // Set Board
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return; // return if canvas is not initialized

    const context = canvas.getContext("2d");
    if (!context) return; // return if context couldn't be obtained

    // adjust canvas size
    canvas.width = 16 * CELL_SIZE;
    canvas.height = 16 * CELL_SIZE;

    // Set Each Tile
    for (let i = 0; i < 256; i++) {
      const byte = i;
      const color = map8bitToRGB(byte);
      const pos = offsetToXY(i);

      // Scale position to account for cell size
      const x = pos.x * CELL_SIZE;
      const y = pos.y * CELL_SIZE;

      context.fillStyle = color; // NOTE:  fill style must be set before drawing
      context.fillRect(x, y, CELL_SIZE, CELL_SIZE);
    }
  });
  return (
    <div style={{ fontFamily: "system-ui, sans-serif", lineHeight: "1.8" }}>
      <h1>Current Board</h1>
      <canvas ref={canvasRef} />
    </div>
  );
}

/**
 *  Maps an 8 Bit Color to RGB
 * 
 *  References : https://en.wikipedia.org/wiki/8-bit_color

 * Mapping
 *
 * Bit    7  6  5  4  3  2  1  0
 *
 * Data   R  R  R  G  G  G  B  B
 * @param byte an 8 Bit Integer
 * @returns RGB CSS String
 */
function map8bitToRGB(byte: number): string {
  // 8 Bit Integer ranging from 0-255
  byte = byte % 256; // In the event this value is greater than 255, SHOULD NOT OCCUR

  console.log(byte);
  const SCALE = 32; // 256 / 8

  // Extract RGB Components
  const RED = (byte >> 5) * SCALE; // Take First 3 Bits by right shifting 5 bits
  const GREEN = ((byte >> 2) & 0x07) * SCALE; // Take Middle 3 Bits by Right Shifting 2 Bits & performing an AND operation against 0x07 | 0b00000111
  const BLUE = (byte & 0x03) * SCALE * 2; // Take Last 2 Bits by performing AND operation between our byte & 0b00000011

  return `rgb(${Math.round(RED)}, ${Math.round(GREEN)}, ${Math.round(BLUE)})`;
}

/**
 * Gets the 2D position of a tile from the provided cell offset
 * @param offset The Cell/Tile Offset
 * @returns The calculated X & Y Position of the tile
 */
function offsetToXY(offset: number) {
  const y = Math.floor(offset / 16);
  const x = offset - 16 * y;
  return { x, y };
}
