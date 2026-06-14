import Image from "next/image";
import Link from "next/link";

interface BlogCardProps {
  image: string;
  date: string;
  readTime: string;
  title: string;
  slug: string;
}

export default function BlogCard({ image, date, readTime, title, slug }: BlogCardProps) {
  return (
    <Link href={`/blog/${slug}`} className="block h-full group">
      <div className="w-full bg-[#FFF5EA] rounded-xl overflow-hidden  transition-all flex flex-col h-full border border-[#F5E6D3]">
        {/* Image Container */}
        <div className="relative w-full aspect-[16/9] bg-gray-200 overflow-hidden">
          <Image src={image} alt={title} fill sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw" className="object-cover transition-transform duration-500" />
        </div>

        {/* Content */}
        <div className="p-5 md:p-6 flex flex-col flex-grow">
          {/* Meta Info */}
          <div className="text-gray-600 text-xs md:text-sm mb-3">
            {date} | {readTime}
          </div>

          {/* Title */}
          <h3 className="font-bold text-base md:text-lg text-black font-sans uppercase leading-snug line-clamp-3 group-hover:text-black/70 transition-colors">
            {title}
          </h3>
        </div>
      </div>
    </Link>
  );
}
