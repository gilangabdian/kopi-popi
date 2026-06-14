import HeroSection from "@/components/public/landing/HeroSection";
import PromoSection from "@/components/public/landing/PromoSection";
import FeaturedProductsSection from "@/components/public/landing/FeaturedProductsSection";
import FeatureSection from "@/components/public/landing/FeatureSection";
import BlogSection from "@/components/public/landing/BlogSection";
import CtaSection from "@/components/public/landing/CtaSection";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center">
      <HeroSection />
      <PromoSection />
      <FeaturedProductsSection />
      <FeatureSection />
      <BlogSection />
      <CtaSection />
    </main>
  );
}
