import { readFileSync, writeFileSync, mkdirSync } from 'fs';

const source = JSON.parse(
  readFileSync('frontend/src/assets/reading-animation-source.json', 'utf8')
);

// #e8c4cf (dark-theme primary, readable on #15100f dark sidebar)
// Normalized RGB: r=232/255=0.9098, g=196/255=0.7686, b=207/255=0.8118
const FILL = [0.9098, 0.7686, 0.8118, 1];

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
