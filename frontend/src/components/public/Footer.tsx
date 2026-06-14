import Image from "next/image";
import Link from "next/link";
import { Icon } from "@iconify-icon/react";

export default function Footer() {
  return (
    <footer className="bg-[#F7DAD9] text-black pt-16 pb-8 px-8 md:px-16 lg:px-16 mt-16 md:mt-32">
      <div className="container mx-auto max-w-[1400px]">
        {/* Main Footer Content */}
        <div className="flex flex-col lg:flex-row justify-between gap-10 lg:gap-8 mb-16">
          {/* Col 1: Logo */}
          <div className="flex-shrink-0">
            <Link href="/" className="block">
              <Image
                src="/logo.png"
                alt="Kopi Popi Logo"
                width={80}
                height={50}
                className="w-[80px] md:w-[80px] h-auto object-contain"
              />
            </Link>
          </div>

          {/* Col 2: Kopi Popi Links */}
          <div className="flex flex-col">
            <h3 className="font-bold text-lg mb-4">Kopi Popi</h3>
            <nav className="flex flex-col space-y-3">
              <Link href="/about" className="text-black/80 hover:text-black text-sm transition-colors">
                About
              </Link>
              <Link href="/lokasi" className="text-black/80 hover:text-black text-sm transition-colors">
                Lokasi
              </Link>
              <Link href="/blog" className="text-black/80 hover:text-black text-sm transition-colors">
                Blog
              </Link>
              <Link href="/membership" className="text-black/80 hover:text-black text-sm transition-colors">
                Membership
              </Link>
            </nav>
          </div>

          {/* Col 3: Bantuan Links */}
          <div className="flex flex-col">
            <h3 className="font-bold text-lg mb-4">Bantuan</h3>
            <nav className="flex flex-col space-y-3">
              <Link href="/contact" className="text-black/80 hover:text-black text-sm transition-colors">
                Hubungi Kami
              </Link>
              <Link href="/terms" className="text-black/80 hover:text-black text-sm transition-colors">
                Syarat dan Ketentuan
              </Link>
              <Link href="/privacy" className="text-black/80 hover:text-black text-sm transition-colors">
                Kebijakan Privasi
              </Link>
              <Link href="/guide" className="text-black/80 hover:text-black text-sm transition-colors">
                Panduan Berbelanja
              </Link>
              <Link href="/faq" className="text-black/80 hover:text-black text-sm transition-colors">
                FAQ
              </Link>
            </nav>
          </div>

          {/* Col 4: Jam Operasional */}
          <div className="flex flex-col">
            <h3 className="font-bold text-lg mb-4">Jam Operasional</h3>
            <div className="flex flex-col space-y-4 text-sm text-black/80">
              <div>
                <p className="font-semibold text-black">Senin - Jumat</p>
                <p>07.00 - 22.00 WIB</p>
              </div>
              <div>
                <p className="font-semibold text-black">Sabtu - Minggu</p>
                <p>08.00 - 23.00 WIB</p>
              </div>
            </div>
          </div>

          {/* Col 5: Social Icons */}
          <div className="flex flex-row gap-3">
            <a
              href="#"
              className="w-10 h-10 bg-white rounded-full flex items-center justify-center hover:bg-gray-100 transition-colors"
              aria-label="Instagram">
              <Icon icon="mdi:instagram" className="text-xl text-black" />
            </a>
            <a
              href="#"
              className="w-10 h-10 bg-white rounded-full flex items-center justify-center hover:bg-gray-100 transition-colors"
              aria-label="Twitter">
              <Icon icon="fa6-brands:x-twitter" className="text-lg text-black" />
            </a>
            <a
              href="#"
              className="w-10 h-10 bg-white rounded-full flex items-center justify-center hover:bg-gray-100 transition-colors"
              aria-label="Facebook">
              <Icon icon="mdi:facebook" className="text-xl text-black" />
            </a>
            <a
              href="#"
              className="w-10 h-10 bg-white rounded-full flex items-center justify-center hover:bg-gray-100 transition-colors"
              aria-label="YouTube">
              <Icon icon="mdi:youtube" className="text-xl text-black" />
            </a>
          </div>
        </div>

        {/* Copyright Bottom */}
        <div className="text-center pt-8">
          <p className="text-black/70 text-xs">@ 2026 Kopi Popi. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
