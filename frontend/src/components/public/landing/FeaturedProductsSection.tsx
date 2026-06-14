"use client";

import { useState, useEffect, useCallback } from "react";
import useEmblaCarousel from "embla-carousel-react";
import { Icon } from "@iconify-icon/react";
import ProductCard from "@/components/ui/ProductCard";

const featuredProducts = [
  {
    image: "/hero3.jpg",
    title: "Kopi Susu Aren",
    description: "Perpaduan espresso pilihan dengan gula aren asli yang manis dan legit.",
    price: 18000,
    badge: "Best Seller",
  },
  {
    image: "/hero3.jpg",
    title: "Caramel Macchiato",
    description: "Espresso klasik dengan sentuhan sirup karamel lembut dan susu creamy.",
    price: 24000,
  },
  {
    image: "/hero3.jpg",
    title: "Matcha Latte",
    description: "Teh hijau matcha premium Jepang berpadu dengan susu segar pilihan.",
    price: 22000,
  },
  {
    image: "/hero3.jpg",
    title: "Americano Ice",
    description: "Kopi hitam pekat dari biji arabica murni dengan sensasi dingin menyegarkan.",
    price: 15000,
  },
  {
    image: "/hero3.jpg",
    title: "Vanilla Latte",
    description: "Kopi susu lembut dengan tambahan aroma ekstrak vanilla yang menenangkan.",
    price: 20000,
  },
];

export default function FeaturedProductsSection() {
  const [isMobile, setIsMobile] = useState(true);

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 768);
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const [emblaRef, emblaApi] = useEmblaCarousel({
    align: "start",
    loop: false,
    dragFree: true,
    watchDrag: isMobile,
  });

  const scrollPrev = useCallback(() => {
    if (emblaApi) emblaApi.scrollPrev();
  }, [emblaApi]);

  const scrollNext = useCallback(() => {
    if (emblaApi) emblaApi.scrollNext();
  }, [emblaApi]);

  return (
    <section className="w-full bg-white py-16 md:py-24 overflow-hidden">
      <div className="container mx-auto px-8 md:px-12 lg:px-20 max-w-[1400px] relative">
        <h2 className="text-3xl md:text-4xl font-sans font-bold text-black mb-8 md:mb-12">Produk Unggulan</h2>

        {/* Carousel Area */}
        <div className="relative group">
          <div className="overflow-hidden" ref={emblaRef}>
            <div className="flex touch-pan-y -ml-4">
              {featuredProducts.map((product, idx) => (
                <div
                  key={idx}
                  className="flex-[0_0_80%] sm:flex-[0_0_50%] md:flex-[0_0_33.333%] lg:flex-[0_0_25%] min-w-0 pl-4">
                  <ProductCard {...product} />
                </div>
              ))}
            </div>
          </div>

          {/* Desktop Arrows */}
          <button
            onClick={scrollPrev}
            className="hidden md:flex absolute -left-12 top-1/2 -translate-y-1/2 w-10 h-10 bg-transparent hover:bg-gray-100 text-black rounded-full items-center justify-center transition-colors"
            aria-label="Previous slide">
            <Icon icon="majesticons:chevron-left-line" width="32" className="text-black font-bold" />
          </button>
          <button
            onClick={scrollNext}
            className="hidden md:flex absolute -right-12 top-1/2 -translate-y-1/2 w-10 h-10 bg-transparent hover:bg-gray-100 text-black rounded-full items-center justify-center transition-colors"
            aria-label="Next slide">
            <Icon icon="majesticons:chevron-right-line" width="32" className="text-black font-bold" />
          </button>
        </div>

        {/* Bottom Button */}
        <div className="mt-12 flex justify-center">
          <button className="bg-[#F7DAD9] hover:bg-[#E4BEBD] text-black font-semibold text-sm md:text-base px-6 py-3 rounded-lg transition-colors">
            Lihat Semua Produk
          </button>
        </div>
      </div>
    </section>
  );
}
