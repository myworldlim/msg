This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.


{
  "name": "chitchat",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "dev:https": "concurrently \"next dev\" \"local-ssl-proxy --source 3001 --target 3000 --cert localhost.pem --key localhost-key.pem\"",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "postinstall": "npm audit fix || true"
  },
  "dependencies": {
    // State Management
    "zustand": "5.0.8",
    
    // API & Data Fetching
    "@tanstack/react-query": "^5.90.2",
    "axios": "^1.7.9",
    "next-auth": "^4.24.5",
    
    // UI & Animation
    "@use-gesture/react": "^10.3.1",
    "framer-motion": "^12.23.22",
    
    // Core
    "next": "15.5.4",
    "react": "19.1.0",
    "react-dom": "19.1.0",
    
    // PWA
    "next-pwa": "5.6.0"
  },
  "devDependencies": {
    // TypeScript
    "@types/node": "20",
    "@types/react": "19",
    "@types/react-dom": "19",
    "@types/next-pwa": "^5.6.9",
    "typescript": "^5",
    
    // Build Tools
    "@babel/core": "^7.28.4",
    
    // Development
    "concurrently": "^9.2.1",
    "local-ssl-proxy": "^2.0.5",
    
    // Linting
    "@eslint/eslintrc": "^3",
    "eslint": "^9",
    "eslint-config-next": "15.2.1"
  }
}