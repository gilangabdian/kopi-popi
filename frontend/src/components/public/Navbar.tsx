"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { motion, AnimatePresence } from "motion/react";
import { Icon } from "@iconify-icon/react";

export default function Navbar() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 w-full bg-white">
      <div className="flex items-center justify-between h-[70px] px-6 md:px-[50px] relative">
        {/* Left Section: Hamburger & Logo */}
        <div className="flex items-center gap-2">
          {/* Mobile Hamburger Icon */}
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="md:hidden flex items-center justify-center p-2 -ml-2 text-black hover:bg-gray-100 rounded-full transition-colors z-10"
            aria-label="Toggle menu">
            <Icon icon={isOpen ? "majesticons:close-line" : "majesticons:menu-line"} width="28" height="28" />
          </button>

          {/* Logo Section */}
          <Link href="/" className="flex items-center hover:opacity-80 transition-opacity">
            <Image
              src="/logo.png"
              alt="Kopi Popi Logo"
              width={50}
              height={20}
              className="w-[35px] md:w-[50px] h-auto object-contain"
              priority
            />
          </Link>
        </div>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center gap-12 text-[15px] font-bold text-black">
          <Link href="/about" className="hover:text-black/60 transition-colors">
            About
          </Link>
          <Link href="/menu" className="hover:text-black/60 transition-colors">
            Menu
          </Link>
          <Link href="/store" className="hover:text-black/60 transition-colors">
            Store
          </Link>
          <Link href="/blog" className="hover:text-black/60 transition-colors">
            Blog
          </Link>
        </nav>

        {/* CTA Buttons */}
        <div className="flex items-center gap-4 md:gap-8 text-sm md:text-[15px] font-bold z-10">
          <Link href="/register" className="hidden md:block text-black hover:text-black/60 transition-colors">
            Daftar
          </Link>
          <Link
            href="/login"
            className="bg-[#F7DAD9] text-black px-5 md:px-8 py-2 md:py-2.5 rounded-full hover:bg-[#E4BEBD] transition-colors">
            Masuk
          </Link>
        </div>
      </div>

      {/* Mobile Menu Dropdown */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.3, ease: "easeInOut" }}
            className="md:hidden overflow-hidden bg-white border-t border-gray-100 shadow-xl absolute w-full">
            <nav className="flex flex-col px-6 py-2 pb-6 text-lg font-bold text-black">
              <Link
                href="/about"
                onClick={() => setIsOpen(false)}
                className="py-4 border-b border-gray-100 hover:text-black/60 transition-colors">
                About
              </Link>
              <Link
                href="/menu"
                onClick={() => setIsOpen(false)}
                className="py-4 border-b border-gray-100 hover:text-black/60 transition-colors">
                Menu
              </Link>
              <Link
                href="/store"
                onClick={() => setIsOpen(false)}
                className="py-4 border-b border-gray-100 hover:text-black/60 transition-colors">
                Store
              </Link>
              <Link
                href="/blog"
                onClick={() => setIsOpen(false)}
                className="py-4 hover:text-black/60 transition-colors">
                Blog
              </Link>
            </nav>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
}
