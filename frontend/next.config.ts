import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",

  // Optimización de tree-shaking para paquetes grandes
  experimental: {
    optimizePackageImports: ['lucide-react', 'date-fns', '@radix-ui/react-dialog', '@radix-ui/react-alert-dialog'],
  },

  images: {
    formats: ['image/avif', 'image/webp'],
    remotePatterns: [
      {
        protocol: "https",
        hostname: "images.unsplash.com",
      },
    ],
  },

  // Headers de caché para assets estáticos
  async headers() {
    return [
      {
        source: '/:all*(svg|jpg|png|webp|avif|woff2)',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
    ];
  },
};

export default nextConfig;

