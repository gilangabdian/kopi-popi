import Image from "next/image";

interface ProductCardProps {
  image: string;
  title: string;
  description: string;
  price: number;
  badge?: string;
}

export default function ProductCard({ image, title, description, price, badge }: ProductCardProps) {
  // Format price to Rupiah
  const formattedPrice = new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price);

  return (
    <div className="w-full bg-[#FFF5EA] rounded-lg overflow-hidden transition-shadow flex flex-col h-full border border-[#F5E6D3]">
      {/* Image Container */}
      <div className="relative w-full aspect-square bg-gray-100">
        <Image src={image} alt={title} fill sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw" className="object-cover" />
        {badge && (
          <div className="absolute top-3 left-3 bg-[#F8D4D4] text-black text-xs font-semibold px-2 py-1 rounded-md shadow-sm">
            {badge}
          </div>
        )}
      </div>

      {/* Content */}
      <div className="p-4 flex flex-col flex-grow">
        <h3 className="font-bold text-lg text-black mb-1 font-sans">{title}</h3>
        <p className="text-gray-600 text-sm mb-4 line-clamp-2 leading-snug flex-grow">{description}</p>

        {/* Bottom Row: Price & Button */}
        <div className="flex items-center justify-between mt-auto pt-2">
          <span className="font-bold text-black">{formattedPrice}</span>
          <button className="bg-[#F7DAD9] hover:bg-[#E4BEBD] text-black text-xs md:text-sm font-semibold px-3 py-1.5 rounded-md transition-colors">
            Beli Sekarang
          </button>
        </div>
      </div>
    </div>
  );
}
