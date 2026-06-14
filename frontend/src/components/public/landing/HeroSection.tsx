import Image from "next/image";

export default function HeroSection() {
  return (
    <section className="w-full bg-white [@media(max-width:400px)]:h-[93vh] [@media(max-width:400px)]:min-h-[450px] h-[95vh] min-h-[450px] md:min-h-0 md:h-[calc(100vh-39px)] flex flex-col overflow-hidden pb-4 md:pb-8">
      {/* Top Text Section */}
      <div className="container mx-auto px-6 text-center max-w-5xl pt-16 md:pt-12 flex-shrink-0">
        <h1 className="[@media(max-width:400px)]:text-5xl [@media(min-width:401px)_and_(max-width:415px)]:text-6xl text-6xl md:text-6xl lg:text-7xl font-serif font-bold text-black leading-[1.15] mb-4">
          Start Your day
          <br />
          with real coffee
        </h1>
        <p className="text-black/80 text-xs md:text-sm max-w-4xl mx-auto leading-relaxed font-regular">
          Kopi Popi hadir untuk menyajikan pengalaman ngopi terbaik. Diracik dari biji kopi Nusantara pilihan yang
          diproses secara presisi, kami memastikan setiap cangkir kopi yang sampai ke tangan Anda mampu membangkitkan
          inspirasi dan menemani setiap langkah hari Anda.
        </p>
      </div>

      {/* Gallery Section with Elliptical Arch Top */}
      <div className="w-full mt-4 flex-1">
        <div
          className="relative w-full h-full flex gap-1 overflow-hidden"
          style={{
            borderTopLeftRadius: "50% 15%",
            borderTopRightRadius: "50% 15%",
          }}>
          {/* Panel 1: Photo (Kecil) */}
          <div className="relative h-full rounded-br-2xl md:rounded-br-[40px] overflow-hidden" style={{ width: "12%" }}>
            <Image
              src="/hero3.jpg"
              alt="Coffee Beans"
              fill
              sizes="(max-width: 768px) 15vw, 12vw"
              className="object-cover"
              priority
            />
          </div>

          {/* Panel 2: Sketch 1 (Lebih besar, tidak di zoom) */}
          <div className="relative h-full rounded-b-2xl md:rounded-b-[40px] overflow-hidden" style={{ width: "27%" }}>
            <Image
              src="/hero2.png"
              alt="Store Sign Sketch"
              fill
              sizes="(max-width: 768px) 30vw, 28vw"
              className="object-cover object-top"
              priority
            />
          </div>

          {/* Panel 3: Text (Tengah) */}
          <div
            className="h-full bg-[#F7DAD9] flex flex-col items-center justify-center p-2 md:p-6 text-center rounded-b-2xl md:rounded-b-[40px] overflow-hidden"
            style={{ width: "24%" }}>
            <h3 className="text-base md:text-md lg:text-[28px] font-serif text-black leading-snug">
              Every cup made to wake you up, warm you up, and keep you moving moving
            </h3>
          </div>

          {/* Panel 4: Sketch 2 (Lebih besar, tidak di zoom) */}
          <div className="relative h-full rounded-b-2xl md:rounded-b-[40px] overflow-hidden" style={{ width: "27%" }}>
            <Image
              src="/hero1.png"
              alt="Coffee Cups Sketch"
              fill
              sizes="(max-width: 768px) 30vw, 28vw"
              className="object-cover object-top"
              priority
            />
          </div>

          {/* Panel 5: Photo 2 (Kecil) */}
          <div className="relative h-full rounded-bl-2xl md:rounded-bl-[40px] overflow-hidden" style={{ width: "12%" }}>
            <Image
              src="/hero4.jpg"
              alt="Coffee Tools"
              fill
              sizes="(max-width: 768px) 15vw, 12vw"
              className="object-cover"
              priority
            />
          </div>
        </div>
      </div>
    </section>
  );
}
