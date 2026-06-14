import Link from "next/link";
import { Icon } from "@iconify-icon/react";

export default function FeatureSection() {
  return (
    <section className="w-full bg-white pt-16 md:pt-16 pb-12">
      <div className="container mx-auto px-6 md:px-12 lg:px-20 max-w-[1400px]">
        {/* Header */}
        <h2 className="text-3xl md:text-4xl font-sans font-bold text-black mb-8 md:mb-12">Featured</h2>
      </div>

      {/* The Grid Container - Full width, no rounded corners */}
      <div className="w-full flex flex-col">
        {/* Top Row: 2 Columns */}
        <div className="flex flex-col md:flex-row w-full">
          {/* Top Left: Order & Pick-up */}
          <div className="w-full md:w-1/2 bg-[#968C83] p-10 md:p-20 flex flex-col items-center text-center">
            {/* Square Illustration Placeholder */}
            <div className="w-40 h-40 bg-gray-200/40 mb-8 flex items-center justify-center">
              <Icon icon="majesticons:shopping-bag-line" className="text-4xl text-white opacity-80" />
            </div>
            <h3 className="text-2xl font-serif text-white mb-4">Pesan Tanpa Antre</h3>
            <p className="text-white/90 text-sm md:text-base max-w-sm mb-6 leading-relaxed">
              Pesan kopi secara online di mana pun kamu berada. Langsung ambil di meja bar tanpa perlu antre saat kamu
              tiba.
            </p>
            <Link href="#" className="text-white/80 text-xs md:text-sm font-semibold underline hover:text-white transition-colors">
              Pelajari Lebih Lanjut
            </Link>
          </div>

          {/* Top Right: Membership */}
          <div className="w-full md:w-1/2 bg-[#F7DAD9] p-10 md:p-20 flex flex-col items-center text-center">
            {/* Square Illustration Placeholder */}
            <div className="w-40 h-40 bg-black/5 mb-8 flex items-center justify-center">
              <Icon icon="majesticons:ticket-line" className="text-4xl text-black opacity-60" />
            </div>
            <h3 className="text-2xl font-serif text-black mb-4">Kopi Popi Club</h3>
            <p className="text-black/70 text-sm md:text-base max-w-sm mb-6 leading-relaxed">
              Daftar jadi member, kumpulkan poin setiap pembelian, dan tukarkan dengan minuman gratis atau kejutan
              spesial.
            </p>
            <Link href="#" className="text-black/70 text-xs md:text-sm font-semibold underline hover:text-black transition-colors">
              Pelajari Lebih Lanjut
            </Link>
          </div>
        </div>

        {/* Bottom Row: 3 Columns */}
        <div className="w-full bg-[#D6D2C4] flex flex-col md:flex-row">
          {/* Bottom 1: Branch */}
          <div className="w-full md:w-1/3 p-10 md:p-16 flex flex-col items-center text-center">
            <h3 className="text-xl font-serif text-black mb-8">Pilih Cabang Terdekat</h3>
            {/* Tilted Illustration Placeholder */}
            <div className="w-40 h-32 bg-black/5 mb-8 -rotate-3 flex items-center justify-center transition-transform hover:rotate-0">
              <Icon icon="majesticons:map-marker-line" className="text-3xl text-black opacity-60" />
            </div>
            <p className="text-black/70 text-sm mb-6 leading-relaxed max-w-xs">
              Gunakan fitur lokasi pintar untuk menemukan cabang Kopi Popi yang paling dekat dengan arah jalanmu.
            </p>
            <Link href="#" className="text-black/70 text-xs font-semibold underline hover:text-black transition-colors mt-auto">
              Pelajari Lebih Lanjut
            </Link>
          </div>

          {/* Bottom 2: Promo */}
          <div className="w-full md:w-1/3 p-10 md:p-16 flex flex-col items-center text-center">
            <h3 className="text-xl font-serif text-black mb-8">Promo Eksklusif</h3>
            {/* Tilted Illustration Placeholder */}
            <div className="w-40 h-32 bg-black/5 mb-8 flex items-center justify-center transition-transform hover:-rotate-2">
              <Icon icon="majesticons:percent-line" className="text-3xl text-black opacity-60" />
            </div>
            <p className="text-black/70 text-sm mb-6 leading-relaxed max-w-xs">
              Dapatkan notifikasi dan klaim berbagai kupon potongan harga yang bisa langsung dipakai saat checkout.
            </p>
            <Link href="#" className="text-black/70 text-xs font-semibold underline hover:text-black transition-colors mt-auto">
              Pelajari Lebih Lanjut
            </Link>
          </div>

          {/* Bottom 3: Live Tracking */}
          <div className="w-full md:w-1/3 p-10 md:p-16 flex flex-col items-center text-center">
            <h3 className="text-xl font-serif text-black mb-8">Pantau Pesanan</h3>
            {/* Tilted Illustration Placeholder */}
            <div className="w-40 h-32 bg-black/5 mb-8 rotate-3 flex items-center justify-center transition-transform hover:rotate-0">
              <Icon icon="majesticons:clock-line" className="text-3xl text-black opacity-60" />
            </div>
            <p className="text-black/70 text-sm mb-6 leading-relaxed max-w-xs">
              Terima update langsung saat pesananmu sedang diproses barista hingga notifikasi pesanan siap untuk
              diambil.
            </p>
            <Link href="#" className="text-black/70 text-xs font-semibold underline hover:text-black transition-colors mt-auto">
              Pelajari Lebih Lanjut
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}
