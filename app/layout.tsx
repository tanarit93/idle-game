import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Idle RPG - Server Authoritative",
  description: "High-performance Golang & Next.js Idle Game",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
