import Image from "next/image";

export default function CtaSection() {
  return (
    <section className="w-full mt-16 md:mt-24">
      <div className="w-full flex flex-col md:flex-row">
        {/* Left Column: Image */}
        <div className="w-full md:w-1/2 relative min-h-[300px] md:min-h-[400px]">
          <Image 
            src="/hero3.jpg" 
            alt="Interior kedai kopi" 
            fill 
            sizes="(max-width: 768px) 100vw, 50vw"
            className="object-cover" 
          />
        </div>

        {/* Right Column: CTA Content */}
        <div className="w-full md:w-1/2 bg-[#FFF5EA] p-12 md:p-24 flex flex-col items-center justify-center text-center">
          <h2 className="text-2xl md:text-3xl font-sans font-medium text-black mb-6 md:mb-8 max-w-md leading-snug">
            Penasaran dengan rasa kopi kami?
          </h2>

          <button className="bg-[#F7DAD9] hover:bg-[#E4BEBD] text-black font-semibold text-sm md:text-base px-6 py-2.5 rounded-lg transition-colors mb-6 md:mb-8">
            Pesan Sekarang
          </button>

          <p className="text-black/70 text-xs md:text-sm max-w-md leading-relaxed">
            Kunjungi cabang terdekat atau pesan langsung via aplikasi. Nikmati racikan biji kopi pilihan dari barista
            terbaik kami, diseduh segar khusus untuk harimu.
          </p>
        </div>
      </div>
    </section>
  );
}
