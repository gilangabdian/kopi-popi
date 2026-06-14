import type { Metadata } from "next";
import "../globals.css";
import { Inter, Playfair_Display } from "next/font/google";
import { cn } from "@/lib/utils";
import Navbar from "@/components/public/Navbar";
import Footer from "@/components/public/Footer";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const playfair = Playfair_Display({ subsets: ["latin"], variable: "--font-playfair" });

export const metadata: Metadata = {
  title: "Kopi Popi | Biji Kopi Nusantara",
  description: "Sistem Point of Sale (POS) dan Inventory Multi-Outlet untuk Kopi-Popi",
  icons: {
    icon: "./logo-web.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id" className={cn("font-sans", inter.variable, playfair.variable)}>
      <body suppressHydrationWarning className="min-h-screen flex flex-col">
        <Navbar />
        <main className="flex-1">{children}</main>
        <Footer />
      </body>
    </html>
  );
}
