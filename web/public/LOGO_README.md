# OpenAuth Logo

## Logo Design

The OpenAuth logo features a modern design that represents:
- **Security**: Shield/keyhole shape representing authentication and security
- **SSO**: Connected nodes representing Single Sign-On functionality
- **Modern**: Gradient colors matching the application theme

## Color Scheme

The logo uses a gradient color scheme:
- Primary: `#6366f1` (Indigo) to `#8b5cf6` (Purple)
- Secondary: `#667eea` (Blue) to `#764ba2` (Purple)

## Files

### SVG Files
- `logo.svg` - Full-size logo (512x512 viewBox)
- `favicon.svg` - Optimized favicon version

### PNG Files (Generated)
- `favicon-16x16.png` - Standard favicon
- `favicon-32x32.png` - Standard favicon
- `favicon-48x48.png` - Standard favicon
- `favicon-64x64.png` - Standard favicon
- `favicon-128x128.png` - Standard favicon
- `favicon-256x256.png` - Standard favicon
- `favicon-512x512.png` - Standard favicon
- `apple-touch-icon.png` - Apple touch icon (180x180)
- `android-chrome-192x192.png` - Android Chrome icon
- `android-chrome-512x512.png` - Android Chrome icon

## Regenerating PNG Files

To regenerate PNG files from the SVG source, run:

```bash
./scripts/generate-favicons.sh
```

**Requirements:**
- ImageMagick (`brew install imagemagick`) or
- Inkscape (`brew install --cask inkscape`)

## Usage

The logo is automatically referenced in:
- `index.html` - All favicon sizes and manifest
- `site.webmanifest` - PWA manifest
- `AppLayout.tsx` - Sidebar logo

## Design Notes

- The logo is designed to be recognizable at small sizes
- Uses high contrast for visibility
- Scalable vector format ensures crisp rendering at any size
- Gradient design matches the application's modern theme
