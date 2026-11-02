// Generate minimal tab bar icons as SVG, then convert to PNG using canvas or manual approach
// We'll create simple SVG files that can be used directly

const fs = require('fs');
const path = require('path');

const dir = path.join(__dirname, 'static', 'tabs');
if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });

// Simple SVG icon generator - uni-app H5 supports SVG
function createSVGIcon(drawFn, color) {
  return `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48">${drawFn(color)}</svg>`;
}

// Ocean waves icon
const oceanIcon = (color) => `
  <path d="M4 28c4-4 8-4 12 0s8 4 12 0 8-4 12 0" fill="none" stroke="${color}" stroke-width="2.5" stroke-linecap="round"/>
  <path d="M4 34c4-4 8-4 12 0s8 4 12 0 8-4 12 0" fill="none" stroke="${color}" stroke-width="2.5" stroke-linecap="round"/>
  <path d="M4 22c4-4 8-4 12 0s8 4 12 0 8-4 12 0" fill="none" stroke="${color}" stroke-width="2.5" stroke-linecap="round"/>
  <circle cx="24" cy="14" r="4" fill="${color}" opacity="0.6"/>
`;

// Bottle icon (for 我的瓶子)
const bottleIcon = (color) => `
  <rect x="18" y="8" width="12" height="4" rx="1" fill="${color}"/>
  <rect x="16" y="12" width="16" height="28" rx="3" fill="none" stroke="${color}" stroke-width="2"/>
  <line x1="20" y1="20" x2="28" y2="20" stroke="${color}" stroke-width="1.5" stroke-linecap="round"/>
  <line x1="20" y1="26" x2="28" y2="26" stroke="${color}" stroke-width="1.5" stroke-linecap="round"/>
  <line x1="20" y1="32" x2="28" y2="32" stroke="${color}" stroke-width="1.5" stroke-linecap="round"/>
`;

// Message/chat icon
const messageIcon = (color) => `
  <rect x="6" y="10" width="36" height="26" rx="4" fill="none" stroke="${color}" stroke-width="2"/>
  <line x1="12" y1="18" x2="30" y2="18" stroke="${color}" stroke-width="2" stroke-linecap="round"/>
  <line x1="12" y1="24" x2="36" y2="24" stroke="${color}" stroke-width="2" stroke-linecap="round"/>
  <line x1="12" y1="30" x2="24" y2="30" stroke="${color}" stroke-width="2" stroke-linecap="round"/>
`;

// User/person icon
const userIcon = (color) => `
  <circle cx="24" cy="16" r="8" fill="none" stroke="${color}" stroke-width="2"/>
  <path d="M10 40c0-8 6-14 14-14s14 6 14 14" fill="none" stroke="${color}" stroke-width="2" stroke-linecap="round"/>
`;

const inactiveColor = '#8899aa';
const activeColor = '#d4a843';

const icons = {
  'home': oceanIcon,
  'home-active': oceanIcon,
  'bottle': bottleIcon,
  'bottle-active': bottleIcon,
  'message': messageIcon,
  'message-active': messageIcon,
  'mine': userIcon,
  'mine-active': userIcon,
};

for (const [name, drawFn] of Object.entries(icons)) {
  const isActive = name.includes('active');
  const color = isActive ? activeColor : inactiveColor;
  const svg = createSVGIcon(drawFn, color);
  fs.writeFileSync(path.join(dir, `${name}.svg`), svg);
}

console.log('SVG icons generated successfully');
