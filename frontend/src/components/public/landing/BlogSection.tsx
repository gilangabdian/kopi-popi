"use client";

import { useState, useEffect, useCallback } from "react";
import useEmblaCarousel from "embla-carousel-react";
import { Icon } from "@iconify-icon/react";
import BlogCard from "@/components/ui/BlogCard";

const latestNews = [
  {
    image: "/hero3.jpg",
    date: "12 Juni 2024",
    readTime: "4 menit",
    title: "RAHASIA DI BALIK ROASTING BIJI KOPI ARABICA TERBAIK KAMI",
    slug: "rahasia-roasting-arabica",
  },
  {
    image: "/hero3.jpg",
    date: "10 Juni 2024",
    readTime: "3 menit",
    title: "5 MANFAAT MINUM KOPI HITAM TANPA GULA DI PAGI HARI",
    slug: "manfaat-kopi-hitam",
  },
  {
    image: "/hero3.jpg",
    date: "5 Juni 2024",
    readTime: "5 menit",
    title: "MENGENAL PERBEDAAN LATTE, CAPPUCCINO, DAN FLAT WHITE",
    slug: "perbedaan-latte-cappuccino",
  },
  {
    image: "/hero3.jpg",
    date: "1 Juni 2024",
    readTime: "4 menit",
    title: "CABANG BARU KOPI POPI KINI HADIR DI KAWASAN SUDIRMAN",
    slug: "cabang-baru-sudirman",
  },
];

export default function BlogSection() {
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
      {/* Container with extra large margins on left and right */}
      <div className="container mx-auto px-8 md:px-16 lg:px-24 max-w-[1400px] relative">
        <h2 className="text-3xl md:text-4xl font-sans font-bold text-black mb-8 md:mb-12">Berita Terbaru</h2>

        {/* Carousel Area */}
        <div className="relative group">
          <div className="overflow-hidden" ref={emblaRef}>
            <div className="flex touch-pan-y -ml-6">
              {latestNews.map((news, idx) => (
                <div key={idx} className="flex-[0_0_85%] sm:flex-[0_0_50%] md:flex-[0_0_33.333%] min-w-0 pl-6">
                  <BlogCard {...news} />
                </div>
              ))}
            </div>
          </div>

          {/* Desktop Arrows (Outside the slider bounds) */}
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
            Lihat Semua Berita
          </button>
        </div>
      </div>
    </section>
  );
}
