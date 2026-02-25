// app/layout.tsx
import type { Metadata } from "next";

import "@/styles/reset.css";
import "./globals.css";

export const metadata: Metadata = {
  title: "ChitChat PWA",
  description: "Scalable chat app on Next.js",
  manifest: '/manifest.json',
};

export const viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  ;

  return (
    <html lang="ru">
      <body>
        <div className="background_main">
              {children}
        </div>
      </body>
    </html>
  );
}