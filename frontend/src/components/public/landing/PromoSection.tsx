"use client";

import { useState, useEffect, useCallback } from "react";
import Image from "next/image";
import { Icon } from "@iconify-icon/react";
import useEmblaCarousel from "embla-carousel-react";
import Autoplay from "embla-carousel-autoplay";

const promoSlides = [
  {
    image: "/hero3.jpg",
    badge: "Weekend Only",
    title: "Weekend Special: Kopi + Pastry",
    description: "Sabtu & Minggu, nikmati paket kopi + pastry pilihan hanya Rp 45.000",
    buttonText: "Klaim Promo",
  },
  {
    image: "/hero1.png",
    badge: "Special Deal",
    title: "Beli 2 Gratis 1",
    description: "Nikmati sore harimu bersama teman-teman dengan promo Beli 2 Gratis 1 untuk semua varian Latte.",
    buttonText: "Lihat Menu",
  },
  {
    image: "/hero4.jpg",
    badge: "New Member",
    title: "Diskon 50% Member Baru",
    description: "Daftar jadi member Kopi Popi sekarang dan dapatkan diskon 50% untuk pesanan pertama Anda.",
    buttonText: "Daftar Sekarang",
  },
  {
    image: "/hero3.jpg",
    badge: "Flash Sale",
    title: "Happy Hour Pagi",
    description: "Setiap jam 07:00 - 09:00 pagi, semua kopi ukuran regular diskon 25% untuk modal semangat harimu.",
    buttonText: "Pesan Sekarang",
  },
];

export default function PromoSection() {
  const [isMobile, setIsMobile] = useState(true);

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 768);
    checkMobile(); // Check initial state
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const [emblaRef, emblaApi] = useEmblaCarousel({ loop: true, dragFree: false, watchDrag: isMobile }, [
    Autoplay({ delay: 8000, stopOnInteraction: false }),
  ]);
  const [currentIndex, setCurrentIndex] = useState(0);

  const scrollPrev = useCallback(() => {
    if (emblaApi) emblaApi.scrollPrev();
  }, [emblaApi]);

  const scrollNext = useCallback(() => {
    if (emblaApi) emblaApi.scrollNext();
  }, [emblaApi]);

  const scrollTo = useCallback(
    (index: number) => {
      if (emblaApi) emblaApi.scrollTo(index);
    },
    [emblaApi],
  );

  useEffect(() => {
    if (!emblaApi) return;

    // Track active index
    const onSelect = () => {
      setCurrentIndex(emblaApi.selectedScrollSnap());
    };

    emblaApi.on("select", onSelect);
    // Initialize first index
    onSelect();

    return () => {
      emblaApi.off("select", onSelect);
    };
  }, [emblaApi]);

  return (
    <section className="w-full bg-[#FFF5EA] py-16 md:py-16 mt-24 overflow-hidden">
      <div className="container mx-auto px-6 max-w-7xl">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 md:gap-12 items-start">
          {/* Left Column: Text & Box */}
          <div className="lg:col-span-4 flex flex-col items-center lg:items-start text-center lg:text-left">
            <h2 className="text-3xl md:text-4xl lg:text-4xl font-sans font-bold text-black mb-8 whitespace-nowrap lg:-mt-10">
              Promo & Penawaran
            </h2>

            {/* The White Box */}
            <div className="bg-white rounded-2xl shadow-sm p-8 md:p-10 w-full max-w-sm flex flex-col items-center text-center">
              <Icon icon="ph:star-four-fill" className="text-[#8B5A2B] w-12 h-12 mb-4" />
              <h3 className="font-serif font-bold text-xl text-black mb-2">Spesial Hari Ini</h3>
              <p className="text-black/70 text-sm leading-relaxed">
                Dapatkan potongan harga khusus untuk setiap pembelian paket Kopi Nusantara Signature kami. Promo
                terbatas!
              </p>
            </div>
          </div>

          {/* Right Column: Carousel */}
          <div className="lg:col-span-8 w-full flex flex-col items-center relative">
            <div className="relative w-full aspect-[4/3] md:aspect-[16/9] rounded-2xl overflow-hidden group bg-white/50">
              {/* Embla Viewport */}
              <div className="overflow-hidden w-full h-full" ref={emblaRef}>
                {/* Embla Container */}
                <div className="flex touch-pan-y w-full h-full">
                  {promoSlides.map((slide, idx) => (
                    <div
                      key={idx}
                      className="flex-[0_0_100%] min-w-0 w-full h-full relative cursor-grab active:cursor-grabbing md:cursor-auto md:active:cursor-auto">
                      <Image src={slide.image} alt={slide.title} fill className="object-cover pointer-events-none" />
                      {/* Dark Gradient Overlay for Text Readability */}
                      <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent flex flex-col justify-end p-6 md:p-10 text-left pointer-events-none">
                        <div className="bg-[#F7DAD9] text-black text-[10px] md:text-sm font-bold px-3 py-1 md:px-4 md:py-1.5 rounded-full w-max mb-2 md:mb-4">
                          {slide.badge}
                        </div>
                        <h3 className="text-white text-xl md:text-4xl font-serif font-bold mb-1 md:mb-3 drop-shadow-md line-clamp-1 md:line-clamp-none">
                          {slide.title}
                        </h3>
                        <p className="text-white/90 text-xs md:text-base mb-4 md:mb-6 max-w-lg drop-shadow-md line-clamp-2 md:line-clamp-none">
                          {slide.description}
                        </p>
                        <button className="bg-[#F7DAD9] hover:bg-[#E4BEBD] text-black font-bold px-4 py-2 md:px-6 md:py-3 rounded-xl w-max transition-colors shadow-lg text-sm md:text-base pointer-events-auto">
                          {slide.buttonText}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Arrows */}
              <button
                onClick={scrollPrev}
                className="hidden md:flex absolute left-4 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/80 hover:bg-white backdrop-blur-sm text-black rounded-full items-center justify-center shadow-lg transition-all duration-300 opacity-0 group-hover:opacity-100 z-10">
                <Icon icon="majesticons:chevron-left-line" width="24" />
              </button>
              <button
                onClick={scrollNext}
                className="hidden md:flex absolute right-4 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/80 hover:bg-white backdrop-blur-sm text-black rounded-full items-center justify-center shadow-lg transition-all duration-300 opacity-0 group-hover:opacity-100 z-10">
                <Icon icon="majesticons:chevron-right-line" width="24" />
              </button>
            </div>

            {/* Dots Navigation */}
            <div className="flex gap-2 mt-6">
              {promoSlides.map((_, idx) => (
                <button
                  key={idx}
                  onClick={() => scrollTo(idx)}
                  className={`h-2 rounded-full transition-all duration-300 ${
                    idx === currentIndex ? "w-8 bg-[#8B5A2B]" : "w-4 bg-gray-300 hover:bg-gray-400"
                  }`}
                  aria-label={`Go to slide ${idx + 1}`}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
