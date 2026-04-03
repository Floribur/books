import { readFileSync, writeFileSync, mkdirSync } from 'fs';

const source = JSON.parse(
  readFileSync('frontend/src/assets/reading-animation-source.json', 'utf8')
);

// #ffffff white — readable on #6d233e brand-red sidebar header
const FILL = [1, 1, 1, 1];

function recolor(obj) {
  if (Array.isArray(obj)) return obj.map(recolor);
  if (obj && typeof obj === 'object') {
    const out = {};
    for (const [k, v] of Object.entries(obj)) {
      // Target solid-color fill: "c" property containing {k: [R, G, B, A]} with 4 numbers
      if (
        k === 'c' &&
        v &&
        typeof v === 'object' &&
        Array.isArray(v.k) &&
        v.k.length === 4 &&
        v.k.every((n) => typeof n === 'number')
      ) {
        out[k] = { ...v, k: FILL };
      } else {
        out[k] = recolor(v);
      }
    }
    return out;
  }
  return obj;
}

mkdirSync('frontend/src/assets', { recursive: true });
const result = recolor(source);
writeFileSync('frontend/src/assets/reading-animation.json', JSON.stringify(result));
console.log('Done: frontend/src/assets/reading-animation.json written with #e8c4cf fill');
